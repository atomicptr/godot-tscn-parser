package godot

type Scene struct {
	Format       int
	ExtResources []string
	SubResources []string
	RootNode     *Node
}
