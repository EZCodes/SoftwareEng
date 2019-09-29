package LCA

import "testing"

func TestLowestCommonAncestor(t *testing.T) {

	tests := []struct {
		tree    Tree
		nodeOne Node
		nodeTwo Node
		want    Node
		wantErr bool
	}{
		{
			// empty tree
			nodeOne: Node{},
			nodeTwo: Node{},
			wantErr: true,
		}, {
			// line tree
			tree: Tree{
				root: Node{
					data: 1,
					childs: []*Node{
						&Node{
							data: 2,
							childs: []*Node{
								&Node{
									data: 3,
									childs: []*Node{
										&Node{
											data: 4,
											childs: []*Node{
												&Node{
													data:   5,
													childs: []*Node{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nodeOne: Node{
				data: 3,
			},
			nodeTwo: Node{
				data: 4,
			},
			wantErr: false,
			want: Node{
				data: 2,
			},
		}, {
			// two sides of the root
			tree: Tree{
				root: Node{
					data: 1,
					childs: []*Node{
						&Node{
							data: 4,
							childs: []*Node{
								&Node{
									data:   5,
									childs: nil,
								},
							},
						},
						&Node{
							data: 2,
							childs: []*Node{
								&Node{
									data:   3,
									childs: nil,
								},
							},
						},
					},
				},
			},
			nodeOne: Node{
				data: 3,
			},
			nodeTwo: Node{
				data: 5,
			},
			wantErr: false,
			want: Node{
				data: 1,
			},
		}, {
			//one side of the root
			tree: Tree{
				root: Node{
					data: 1,
					childs: []*Node{
						&Node{
							data:   4,
							childs: nil,
						},
						&Node{
							data: 2,
							childs: []*Node{
								&Node{
									data: 3,
									childs: []*Node{
										&Node{
											data:   5,
											childs: nil,
										},
									},
								},
								&Node{
									data:   6,
									childs: nil,
								},
							},
						},
					},
				},
			},
			nodeOne: Node{
				data: 5,
			},
			nodeTwo: Node{
				data: 6,
			},
			wantErr: false,
			want: Node{
				data: 2,
			},
		}, {
			// more than 2 childs
			tree: Tree{
				root: Node{
					data: 1,
					childs: []*Node{
						&Node{
							data: 2,
							childs: []*Node{
								&Node{
									data:   3,
									childs: nil,
								},
								&Node{
									data:   4,
									childs: nil,
								},
								&Node{
									data:   5,
									childs: nil,
								},
							},
						},
					},
				},
			},
			nodeOne: Node{
				data: 3,
			},
			nodeTwo: Node{
				data: 4,
			},
			wantErr: false,
			want: Node{
				data: 2,
			},
		}, {
			//node is not in the tree
			tree: Tree{
				root: Node{
					data: 1,
					childs: []*Node{
						&Node{
							data: 4,
							childs: []*Node{
								&Node{
									data:   5,
									childs: nil,
								},
							},
						},
						&Node{
							data: 2,
							childs: []*Node{
								&Node{
									data:   3,
									childs: nil,
								},
							},
						},
					},
				},
			},
			nodeOne: Node{
				data: 3,
			},
			nodeTwo: Node{
				data: 7,
			},
			wantErr: false,

			want: Node{
				data: 0,
			},
		},
	}
	for index, test := range tests {
		t.Run(string(index), func(t *testing.T) {
			fNode, err := test.tree.LowestCommonAncestor(&(test.nodeOne), &(test.nodeTwo))
			if !(test.wantErr) && err != nil {
				t.Errorf("Error has occured: %q", err)
			}
			if !(test.wantErr) && test.want.data != fNode.data {
				t.Errorf("Found ancestor is incorrect. Want: %d, Got: %d", test.want.data, fNode.data)
			}
		})
	}
}
