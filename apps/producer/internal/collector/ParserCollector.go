package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"producer/internal/kafka-producer"
	"producer/internal/parsers/interfaces"
	discogs_functions2 "producer/internal/parsers/parser-discogs/discogs-functions"
	discogs_structs2 "producer/internal/parsers/parser-discogs/discogs-structs"
	parser_lastfm "producer/internal/parsers/parser-lastfm/lastfm-functions"
	"time"
)

func ParserCollectorAlbum(lastfmAlbum, discogsAlbum interfaces.IAlbum) []byte {
	// log.Printf("Length 1 : %v, Length 2 : %v\n", lastfmAlbum.GetTracksLen(), discogsAlbum.GetTracksLen())

	//log.Println("Artist id before hash :", discogsAlbum.GetArtistsId())
	album := AlbumDBBuilder(lastfmAlbum, discogsAlbum) // building db structure object from two elements
	//log.Println("Artist id after hash :", album.ArtistHash)
	data, err := json.Marshal(album)
	if err != nil {
		log.Println("Error marshalling album :", discogsAlbum.GetTitle())
	}
	data = bytes.Trim(data, "\x00")

	return data
}

func ParserCollectorTrack(lastfmAlbum, discogsAlbum interfaces.IAlbum) [][]byte {
	var allTracksData [][]byte // creating list of tracks for later usage in consumer
	var albumTracks []interfaces.ITrack
	if lastfmAlbum.GetTracksLen() >= discogsAlbum.GetTracksLen() { // taking the album with the biggest length
		albumTracks = lastfmAlbum.GetTracks()
	} else {
		albumTracks = discogsAlbum.GetTracks()
	}

	for i := 0; i < len(albumTracks); i++ { // iterating over the album tracks list to convert them to the db table format
		//	log.Printf("---%v. Original value : %v", i+1, albumTracks[i])

		trackTmp := TrackDBBuilder(albumTracks[i], discogsAlbum)
		log.Println("--Artist id according to track :", trackTmp.ArtistHash)
		log.Println("---Album id according to track :", trackTmp.AlbumHash)
		//	log.Printf("Type : %T", trackTmp)

		//log.Printf("---%v. Final track ArtistHash : %v\n", i+1, trackTmp.ArtistHash)
		//log.Printf("---%v. Final track AlbumHash : %v\n", i+1, trackTmp.AlbumHash)

		data, err := json.Marshal(trackTmp)
		if err != nil {
			log.Println("Error marshalling track :", albumTracks[i].GetTitle())
		}
		data = bytes.Trim(data, "\x00")
		allTracksData = append(allTracksData, data)
	}

	return allTracksData
}

func CheckErrorParser(err error, message string) {
	if err != nil {
		log.Println(message)
	}
}

func ParserCollectorArtistWithReleases(ArtistName string) {
	// lastfm Artist data got
	lastfmArtist := parser_lastfm.ReadArtist(ArtistName) // artist of lastfm in struct

	// discogs Artist data got
	discogsArtistID := discogs_functions2.CreateRequestByName(ArtistName, "artist")
	discogsArtistData := discogs_functions2.ReadArtistById(discogsArtistID)
	var discogsArtist discogs_structs2.DiscogsArtistJson // artist of discogs in struct
	err := json.Unmarshal(discogsArtistData, &discogsArtist)
	if err != nil {
		log.Println("No such artist as : ", ArtistName)
		return
	}
	// all pages from artistId from discogs got
	// we will iterate over these pages, searching for releases which are marked by tag "master" - this tag is a sign that
	// this release is a originally created album, not a custom compilation etc
	discogsReleasesData := discogs_functions2.ReadReleasesByArtistId(discogsArtistID) // form discogs we get the releases list
	var discogsReleasesPages discogs_structs2.DiscogsPagesReleases                    // all releases by artist id
	err = json.Unmarshal(discogsReleasesData, &discogsReleasesPages)
	CheckErrorParser(err, fmt.Sprintln("Error unmarshalling album data for artist ", ArtistName))

	albumsNum := 0
	for _, r := range discogsReleasesPages.Releases { // iterating over artist releases
		if r.Artist == ArtistName && r.Type == "master" { // filter only on master releases
			albumsNum++
			// here we got the right name for the album, starting working with it....
			time.Sleep(30 * time.Second)                                // anti antiDDOS pause
			discogsAlbumData := discogs_functions2.ReadMasterById(r.ID) // reading album data by id
			// BUT sometimes id leads to another album that has nothing to do
			// with the artist or searched album !!!! (discogs bug, I guess...)
			// fixing it by another function, that builds AlbumDB from discogs.releases + lastfm.album

			//log.Printf("\n%v. ID %v, Title : %v, Type : %v\n", i+1, r.ID, r.Title, r.Type) // writing info into log
			log.Printf("Url for release : %s\n", r.ResourceURL)

			var masterAlbum discogs_structs2.DiscogsMasters
			err := json.Unmarshal(discogsAlbumData, &masterAlbum) // unmarshalling data got by albumId
			if err != nil {
				log.Println("Error unmarshalling data of album by number of : ", albumsNum+1)
				continue
			}
			lastfmAlbum := parser_lastfm.ReadAlbum(ArtistName, r.Title)

			if lastfmAlbum.GetTitle() == "" {
				log.Println("Album does not exist in lastfm")
				continue
			}

			var aDb []byte
			var tDb [][]byte

			//v := masterAlbum.ResourceURL == r.ResourceURL
			//log.Println("Discogs Master Resourses Url Equals Discgos Release Resources Url : ", v)

			if masterAlbum.Title == r.Title { // checking the albumTmp leads to the same album that was requested from releases list
				aDb = ParserCollectorAlbum(lastfmAlbum, masterAlbum) // full album thing
				tDb = ParserCollectorTrack(lastfmAlbum, masterAlbum)
				//log.Println("////////////Artist id according to album : ", masterAlbum.GetArtistsId())

			} else {
				aDb = ParserCollectorAlbum(lastfmAlbum, r) // if id leads to the mistake
				tDb = ParserCollectorTrack(lastfmAlbum, r)
				//log.Println("////////////Artist id according to album : ", r.GetArtistsId())
			}

			var checkAlbumDB AlbumDB
			err = json.Unmarshal(aDb, &checkAlbumDB) // unmarshalling data just to print it in log (you may skip it for time saving)
			if err != nil {
				log.Println("Error unmarshalling data (after parsing) of album by number of : ", albumsNum+1)
				continue
			}

			//log.Printf("%v. Final title : %v, Final tracks number : %v\n",
			//	i+1, checkAlbumDB.Name, checkAlbumDB.TrackCount) // writing unmarshalled info into log
			//log.Printf("%v. Final album Artisthash : %v\n", i+1, checkAlbumDB.ArtistHash)
			//log.Printf("%v. Final album AlbumHash : %v\n", i+1, checkAlbumDB.AlbumHash)
			aDb = bytes.Trim(aDb, "\x00") // trimming to get rid of useless bytes, that cause problems in the consumer

			kafka_producer.Produce("Album", aDb, nil)
			// log.Println("Album sent...")

			timeout := time.After(10 * time.Second)
			chanTopic := make(chan string)
			for i := 0; i < len(tDb); i++ {
				go kafka_producer.Produce("Track", tDb[i], chanTopic)
			}
			for i := 0; i < len(tDb); i++ {
				select {
				case c := <-chanTopic:
					log.Printf("Value sent to kafka topic : %v\n", c)
				case <-timeout:
					log.Println("Error sending value to kafka topic (timeout)!")
				}
			}
		}
	}

	artist := ArtistDBBuilder(lastfmArtist, discogsArtist) // creating artist table element from two objects

	dataArtist, err := json.Marshal(artist)
	if err != nil {
		log.Println("Error marshalling artist :", ArtistName)
	}
	dataArtist = bytes.Trim(dataArtist, "\x00")
	kafka_producer.Produce("Artist", dataArtist, nil)
}
