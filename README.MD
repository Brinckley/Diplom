First we need to start kafka:  
by this comman if it is out firts time starting app:  
`sudo docker compose up --build init-kafka`  
of by this for second and later:  
`sudo docker compose up init-kafka`  
The topics were created.  
  
Now let's start the producer.  
`sudo docker compose up --no-deps --build producer`  
This will put some info in kafka (one artist, one album, one track yet).  
  
Same thing with consumer.  
`sudo docker compose up --no-deps --build consumer`  
However, we may need the db to start also. Yet consumer depend only on db, so it will start automatically,
but it may be better to start it in solo.  