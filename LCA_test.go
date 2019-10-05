package LCA

import "testing"

func TestLowestCommonAncestor(t *testing.T) {

	tests := []struct {
		graph   Graph
		vOne    *Vertex
		vTwo    *Vertex
		want    Vertex
		wantErr bool
	}{
		{
			// empty tree
			vOne:    nil,
			vTwo:    nil,
			wantErr: true,
		}, {
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 2},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 2},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 3},
					},
					&Edge{
						dest: &Vertex{data: 5},
						src:  &Vertex{data: 2},
					},
					&Edge{
						dest: &Vertex{data: 6},
						src:  &Vertex{data: 3},
					},
					&Edge{
						dest: &Vertex{data: 7},
						src:  &Vertex{data: 5},
					},
					&Edge{
						dest: &Vertex{data: 7},
						src:  &Vertex{data: 4},
					},
				},
			},
			vOne: &Vertex{
				data: 7,
			},
			vTwo: &Vertex{
				data: 6,
			},
			wantErr: false,
			want: Vertex{
				data: 3,
			},
		}, {
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 2},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 2},
					},
					&Edge{
						dest: &Vertex{data: 5},
						src:  &Vertex{data: 3},
					},
					&Edge{
						dest: &Vertex{data: 6},
						src:  &Vertex{data: 3},
					},
					&Edge{
						dest: &Vertex{data: 6},
						src:  &Vertex{data: 2},
					},
				},
			},
			vOne: &Vertex{
				data: 6,
			},
			vTwo: &Vertex{
				data: 4,
			},
			wantErr: false,
			want: Vertex{
				data: 2,
			},
		}, {
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 5},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 2},
					},
					&Edge{
						dest: &Vertex{data: 5},
						src:  &Vertex{data: 2},
					},
				},
			},
			vOne: &Vertex{
				data: 5,
			},
			vTwo: &Vertex{
				data: 3,
			},
			wantErr: false,
			want: Vertex{
				data: 1,
			},
		}, {
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 2},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 2},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 3},
					},
					&Edge{
						dest: &Vertex{data: 5},
						src:  &Vertex{data: 4},
					},
				},
			},
			vOne: &Vertex{
				data: 4,
			},
			vTwo: &Vertex{
				data: 3,
			},
			wantErr: false,
			want: Vertex{
				data: 2,
			},
		}, {
			//no common ancestor
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 4},
						src:  &Vertex{data: 2},
					},
				},
			},
			vOne: &Vertex{
				data: 2,
			},
			vTwo: &Vertex{
				data: 3,
			},
			wantErr: false,

			want: Vertex{
				data: 0,
			},
		}, {
			//node not in the graph
			graph: Graph{
				edges: []*Edge{
					&Edge{
						dest: &Vertex{data: 3},
						src:  &Vertex{data: 1},
					},
					&Edge{
						dest: &Vertex{data: 2},
						src:  &Vertex{data: 1},
					},
				},
			},
			vOne: &Vertex{
				data: 4,
			},
			vTwo: &Vertex{
				data: 2,
			},
			wantErr: false,

			want: Vertex{
				data: 0,
			},
		},
	}
	for index, test := range tests {
		t.Run(string(index), func(t *testing.T) {
			fNode, err := test.graph.LowestCommonAncestor(test.vOne, test.vTwo)
			if !(test.wantErr) && err != nil {
				t.Errorf("Error has occured: %q", err)
			}
			if !(test.wantErr) && test.want.data != fNode.data {
				t.Errorf("Found ancestor is incorrect. Want: %d, Got: %d", test.want.data, fNode.data)
			}
		})
	}
}
