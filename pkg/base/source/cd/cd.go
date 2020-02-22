package cd

const table = "CD"

type columns struct {
	cdID  int    `name:"cd_id"`
	title string `name:"title"`
	year  int    `name:"year"`
}

type relationships struct {
	tracks relationship.HasMany{
		fromTable: table,
		
	artists []Artist
}

type CD struct {
	columns columns
	relationships
}

func (c *CD) CDID() int {
	return c.columns.cdID
}

func (c *CD) Title() string {
	return c.columns.title
}

func (c *CD) Year() int {
	return c.columns.year
}

func (c *CD) Artists() []Artist {
	if c.artists == nil {
		c.loadArtists()
	}
	return c.relationships.artists
}

func (c *CD) loadArtists() {

}
