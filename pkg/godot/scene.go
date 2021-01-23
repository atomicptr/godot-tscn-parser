package godot

type Scene struct {
	Format       int
	ExtResources map[int64]*ExtResource
	SubResources map[int64]*SubResource
	RootNode     *Node
	MetaData
}
