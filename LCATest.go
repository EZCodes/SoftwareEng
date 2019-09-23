package LCA

import "testing"


(t *testing.T) testFindPath{


}

(t *testing.T) testLowestCommonAncestor{

	tests :=	[]struct {
		tree Tree
		Node nodeOne
		Node nodeTwo
		want Node
		wantErr bool
	}{
	{
	// empty tree
		nodeOne : nil,
		nodeTwo : nil,
		wantErr : true,
	},{
	// line tree
		tree : Tree{ 
			root: Node {
				data: 1,
				childs : []*Node {
					{
						&Node{
							data: 2,
							childs : []*Node {
								{
									&Node{
										data: 3,
										childs : []*Node {
											{
												&Node{
													data: 4,
													childs : []*Node {
													{
														&Node{
															data: 5,
															childs : nil
															},
														}
													},
												},
											}
										},
									},
								}						
							},
						},
					}
				},
			},
		},
		nodeOne: Node{
			data: 3,		
		},
		nodeTwo: Node{
			data: 4,
		} 
		wantErr: false,
	},{
	// two sides of the root
	tree : Tree{ 
			root: Node {
				data: 1,
				childs : []*Node {
					{
						&Node{
							data: 4,
							childs : []*Node {
								{
									&Node{
										data: 5,
										childs : nil
										},
									},
								}						
							},					
						},
						&Node{
							data: 2,
							childs : []*Node {
								{
									&Node{
										data: 3,
										childs : nil
										},
									},
								}						
							},
						},
					}
				},
			},
		},
		nodeOne: Node{
			data: 3,		
		},
		nodeTwo: Node{
			data: 5,
		} 
		wantErr: false,	
	},{
	//one side of the root
	tree : Tree{ 
		root: Node {
			data: 1,
			childs : []*Node {
				{
					&Node{
						data: 4,
						childs : nil 				
						},					
					},
					&Node{
						data: 2,
						childs : []*Node {
							{
								&Node{
									data: 3,
									childs : []*Node {
										{
											&Node{
												data: 5,
												childs : nil
											},
										}
									},
								},
								&Node{
									data: 6,
									childs : nil
								},
							}						
						},
					},
				}
			},
		},
	},
	nodeOne: Node{
		data: 5,		
	},
	nodeTwo: Node{
		data: 6,
	} 
	wantErr: false,	
	},{
	// more than 2 childs
	tree : Tree{ 
			root: Node {
				data: 1,
				childs : []*Node {
					{
						&Node{
							data: 2,
							childs : []*Node {
								{
									&Node{
										data: 3,
										childs : nil,
										},
									&Node{
										data: 4,
										childs : nil,
									},
									&Node{
										data: 5,
										childs : nil,
									},					
								}						
							},										
						},
					}
				},
			},
		},
		nodeOne: Node{
			data: 3,		
		},
		nodeTwo: Node{
			data: 4,
		} 
		wantErr: false,	
	
	
	},{
	//node is not in the tree
	tree : Tree{ 
			root: Node {
				data: 1,
				childs : []*Node {
					{
						&Node{
							data: 4,
							childs : []*Node {
								{
									&Node{
										data: 5,
										childs : nil
										},
									},
								}						
							},					
						},
						&Node{
							data: 2,
							childs : []*Node {
								{
									&Node{
										data: 3,
										childs : nil
										},
									},
								}						
							},
						},
					}
				},
			},
		},
		nodeOne: Node{
			data: 3,		
		},
		nodeTwo: Node{
			data: 7,
		} 
		wantErr: false,	
	
	
	},{
	// 3 descendants TODO
	
	},
	}

}

