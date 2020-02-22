package eg

import "fmt"

func main() {
	s := myschema.Connect(db)

	album := s.Album().Find(14)
	album.SetTitle("Physical Graffiti")
	title := album.Title()
	err := album.Update()
	album.DiscardChanges()

	newAlbum := s.Album().New()
	newAlbum.SetTitle("Wish You Were Here")
	newAlbum.SetArtist("Pink Floyd")
	err = newAlbum.Insert()
	err = newAlbum.Delete()

	s.Album().Search(s.Album().Where.ArtistEQ("Falco")).Delete()

	rs := s.Album().Search(s.Album().Where.ArtistEQ("Santana"))
	for rs.Next() {
		album := rs.Album()
		fmt.Printf("%s - %s\n", album.Artist(), album.Title())
		album.SetYear(2001)
		err = album.Update()
	}
	if err = rs.Error(); err != nil {
		// do something
	}

	err = rs.Update(s.AlbumUpdate.Year("2001"))

	album, err = s.Album().Search(s.Album().Where.Artist.EQ("Santana")).All()
	for _, album := range album {
		fmt.Printf("%s - %s\n", album.Artist(), album.Title())
	}

	album, err := s.Album().Search(s.Album().Where.Artist.EQ("Santana")).One()

	rs = s.Album().Search(s.Album().Where.Artist.Like("Jimi%"))
	rs = s.Album().Search(s.Literal("artist = $1 AND year = $2", "Peter Frampton", 1986))

	rs = s.Album().Search(
		s.Album().Where.Artist.NE("Janis Joplin"),
		s.Album().Where.Year.LT(1980),
		s.Album().Where.AlbumID.In(1, 14, 15, 65, 43),
	)

	rs = s.Album().Search(s.Album().Where.Artist.EQ("Bob Marley")).
		Limit(2).OrderBy(s.Album().OrderBy.Year(s.OrderByDESC))
}

func joins() {
	rs = s.Album().Search(
		s.Album().Join.Artist(),
		s.Artist().Where.Name.EQ("Bob Marley"),
	).Prefetch(s.Album().Join.Artist)
}
