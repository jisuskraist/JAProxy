package balance

//RouteMapping represents a mapping of a domain with it's targets destinations.
//There could be one or more targets.
type RouteMapping struct {
	Domain  string
	Targets []string
}
