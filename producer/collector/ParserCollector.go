package collector

import (
	"bytes"
	"encoding/json"
	"log"
	"producer/parsers/interfaces"
	"producer/parsers/parser-discogs/discogs-functions"
	"producer/parsers/parser-discogs/discogs-structs"
	parser_lastfm "producer/parsers/parser-lastfm/lastfm-functions"
	"time"
)

func ParserCollectorArtist(ArtistName string) []byte {
	lastfmArtist := parser_lastfm.ReadArtist(ArtistName) // artist of lastfm in struct
	discogsArtistID := discogs_functions.CreateRequestByName(ArtistName, "artist")
	discogsArtistData := discogs_functions.ReadArtistById(discogsArtistID)
	var discogsArtist discogs_structs.DiscogsArtistJson // artist of discogs in struct
	json.Unmarshal(discogsArtistData, &discogsArtist)

	artist := ArtistDBBuilder(lastfmArtist, discogsArtist, 2) // building db structure object from two elements
	data, err := json.Marshal(artist)
	if err != nil {
		log.Println("Error marshalling artist :", ArtistName)
	}
	data = bytes.Trim(data, "\x00")

	return data
}

// aName      string
// release    time.Time
// urlLastfm  string
// urlDicogs  string
// picture    string
// trackCount int
func ParserCollectorAlbum(lastfmAlbum, discogsAlbum interfaces.IAlbum) []byte {
	log.Printf("Length 1 : %v, Length 2 : %v\n", lastfmAlbum.GetTracksLen(), discogsAlbum.GetTracksLen())

	album := AlbumDBBuilder(lastfmAlbum, discogsAlbum) // building db structure object from two elements
	data, err := json.Marshal(album)
	if err != nil {
		log.Println("Error marshalling album :", discogsAlbum.GetTitle())
	}
	data = bytes.Trim(data, "\x00")

	return data
}

// tName     string
// urlLastfm string
// release   time.Time
// lyrics    string
func ParserCollectorTrack(lastfmAlbum, discogsAlbum interfaces.IAlbum) [][]byte {
	var allTracksData [][]byte // creating list of tracks for later usage in consumer
	var albumTracks []interfaces.ITrack
	if lastfmAlbum.GetTracksLen() > discogsAlbum.GetTracksLen() { // taking the album with the biggest length
		albumTracks = lastfmAlbum.GetTracks()
	} else {
		albumTracks = discogsAlbum.GetTracks()
	}

	for _, t := range albumTracks { // iterating over the album tracks list to convert them to the db table format
		trackTmp := TrackDBBuilder(t)
		data, err := json.Marshal(trackTmp)
		if err != nil {
			log.Println("Error marshalling track :", t.GetTitle())
		}
		data = bytes.Trim(data, "\x00")
		allTracksData = append(allTracksData, data)
	}

	return allTracksData
}

func ParserCollectorArtistWithReleases(ArtistName string) ([]byte, [][]byte, [][][]byte) {
	lastfmArtist := parser_lastfm.ReadArtist(ArtistName) // artist of lastfm in struct

	discogsArtistID := discogs_functions.CreateRequestByName(ArtistName, "artist")
	discogsArtistData := discogs_functions.ReadArtistById(discogsArtistID)
	var discogsArtist discogs_structs.DiscogsArtistJson // artist of discogs in struct
	json.Unmarshal(discogsArtistData, &discogsArtist)

	discogsReleasesData := discogs_functions.ReadReleasesByArtistId(discogsArtistID) // form discogs we get the releases list
	var discogsReleasesPages discogs_structs.DiscogsPagesReleases                    // all releases by artist id
	json.Unmarshal(discogsReleasesData, &discogsReleasesPages)
	log.Println(len(discogsReleasesPages.Releases)) // printing number of ALL releases by selected artist

	var albumsData [][]byte                           // all releases data stored here
	var tracksData [][][]byte                         // all tracks data stored here
	for i, r := range discogsReleasesPages.Releases { // iterating over artist releases
		if r.Artist == ArtistName && r.Type == "master" {
			time.Sleep(15 * time.Second)                               // anti antiDDOS pause
			discogsAlbumData := discogs_functions.ReadMasterById(r.ID) // reading album data by id
			// BUT sometimes id leads to another album that has nothing to do
			// with the artist or searched album !!!!
			// fixing it by another function, that builds AlbumDB from discogs.releases + lastfm.album

			log.Printf("\n%v. ID %v, Title : %v, Type : %v\n", i+1, r.ID, r.Title, r.Type) // writing info into log
			log.Println(r.ResourceURL)
			var masterAlbum discogs_structs.DiscogsMasters
			json.Unmarshal(discogsAlbumData, &masterAlbum) // unmarshalling data got by albumId

			lastfmAlbum := parser_lastfm.ReadAlbum(ArtistName, r.Title)

			var aDb []byte
			var tDb [][]byte
			if masterAlbum.Title == r.Title { // checking the albumTmp leads to the same album that was requested from releases list
				aDb = ParserCollectorAlbum(lastfmAlbum, masterAlbum) // full album thing
				tDb = ParserCollectorTrack(lastfmAlbum, masterAlbum)
			} else {
				aDb = ParserCollectorAlbum(lastfmAlbum, r) // if id leads to the mistake
				tDb = ParserCollectorTrack(lastfmAlbum, r)
			}

			var checkAlbumDB AlbumDB
			json.Unmarshal(aDb, &checkAlbumDB)                                                                               // unmarshalling data just to print it in log (you may skip it for time saving)
			log.Printf("%v. Final title : %v, Final tracks number : %v\n", i+1, checkAlbumDB.AName, checkAlbumDB.TrackCount) // writing unmarshalled info into log

			aDb = bytes.Trim(aDb, "\x00")        // trimming to get rid of useless bytes, that cause problems in the consumer
			albumsData = append(albumsData, aDb) // adding album data to the list

			tracksData = append(tracksData, tDb) // adding track data to the list
		}
	}

	//log.Printf("Releases amount for artist %v is %v", ArtistName, len(discogsReleasesPages.Releases))

	artist := ArtistDBBuilder(lastfmArtist, discogsArtist, len(albumsData)) // creating artist table element from two objects
	dataArtist, err := json.Marshal(artist)
	if err != nil {
		log.Println("Error marshalling artist :", ArtistName)
	}
	dataArtist = bytes.Trim(dataArtist, "\x00")

	return dataArtist, albumsData, tracksData
}
