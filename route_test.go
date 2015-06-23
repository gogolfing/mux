package mux

import (
	"reflect"
	"testing"
)

var (
	empty    = newRoute("")
	a        = newRoute("a")
	b        = newRoute("b")
	hello    = newRoute("hello")
	another  = newRoute("another")
	boy      = newRoute("boy")
	chooses  = newRoute("chooses")
	division = newRoute("division")
	elephant = newRoute("elephant")
	frogs    = newRoute("frogs")
	giraffe  = newRoute("giraffe")
)

func TestNewRoute(t *testing.T) {
	tests := []struct {
		path     string
		children []*Route
	}{
		{"path", nil},
		{"path", []*Route{}},
		{"another path", []*Route{&Route{}, &Route{}}},
	}
	for _, test := range tests {
		route := newRoute(test.path, test.children...)
		if route.path != test.path {
			t.Errorf("route.path = %v want %v", route.path, test.path)
		}
		if !areRoutesEqual(route.children, test.children) {
			t.Fail()
		}
	}
}

func TestRoute_insertSplitChild(t *testing.T) {
	tests := []struct {
		root    *Route
		oldPath string
		newPath string
	}{
		{
			newRoute("",
				newRoute("hello"),
			),
			"hello",
			"he",
		},
		{
			newRoute("",
				newRoute("another"),
				newRoute("baseball"),
				newRoute("car"),
				newRoute("diamond"),
				newRoute("hello"),
				newRoute("world"),
			),
			"hello",
			"hell",
		},
	}
	for _, test := range tests {
		expectedNumberChildren := len(test.root.children)
		_, expectedOldRoute, remainingPath := test.root.findSubRoute(test.oldPath)
		_, expectedIndex, _ := test.root.findChildWithCommonPrefix(test.oldPath)
		if expectedOldRoute == nil || len(remainingPath) > 0 {
			t.Errorf("route.findSubRoute(%q) = %v, %q is incorrect", test.oldPath, expectedOldRoute, remainingPath)
		}
		expectedOldRoute.routeHandler = &routeHandler{}

		expectedNewRoute := test.root.insertSplitChild(test.newPath)

		_, oldRoute, remaingPath := test.root.findSubRoute(test.oldPath)
		if oldRoute != expectedOldRoute || len(remaingPath) > 0 || !reflect.DeepEqual(oldRoute, expectedOldRoute) {
			t.Errorf("route.insertSplitChild(%q) oldRoute = %v want %v", test.oldPath, oldRoute, expectedOldRoute)
		}
		_, newRoute, remaingPath := test.root.findSubRoute(test.newPath)
		if newRoute != expectedNewRoute || newRoute.path != test.newPath || len(remainingPath) > 0 {
			t.Errorf("route.insertSplitChild(%q) newRoute = %v want %v", test.newPath, newRoute, expectedNewRoute)
		}
		_, index, _ := test.root.findChildWithCommonPrefix(test.newPath)
		if index != expectedIndex {
			t.Errorf("route.insertSplitChild(%q) index = %v want %v", test.newPath, index, expectedIndex)
		}
		numberChildren := len(test.root.children)
		if expectedNumberChildren != numberChildren {
			t.Errorf("route.insertSplitChild(%q) number of children = %v want %v", test.newPath, numberChildren, expectedNumberChildren)
		}
	}
}

func TestRoute_findSubRoute(t *testing.T) {
	t.Log("root with no children")
	root := newRoute("")
	testRoute_findSubRoute(t, root, "", root, nil, "")
	testRoute_findSubRoute(t, root, "hello", root, nil, "hello")

	t.Log("root with path and no children")
	root = newRoute("root")
	testRoute_findSubRoute(t, root, "", root, nil, "")
	testRoute_findSubRoute(t, root, "hello", root, nil, "hello")

	t.Log("root with single child")
	helloRoute := newRoute("hello")
	root = newRoute("", helloRoute)
	testRoute_findSubRoute(t, root, "", root, nil, "")
	testRoute_findSubRoute(t, root, "he", root, helloRoute, "he")
	testRoute_findSubRoute(t, root, "hello", root, helloRoute, "")
	testRoute_findSubRoute(t, root, "hello, world", helloRoute, nil, ", world")
	testRoute_findSubRoute(t, root, "another", root, nil, "another")

	t.Log("root with single child and single grandchild")
	worldRoute := newRoute(", world")
	helloRoute = newRoute("hello", worldRoute)
	root = newRoute("", helloRoute)
	testRoute_findSubRoute(t, root, "", root, nil, "")
	testRoute_findSubRoute(t, root, "he", root, helloRoute, "he")
	testRoute_findSubRoute(t, root, "hello", root, helloRoute, "")
	testRoute_findSubRoute(t, root, "hello, world", helloRoute, worldRoute, "")
	testRoute_findSubRoute(t, root, "another", root, nil, "another")
	testRoute_findSubRoute(t, root, "hello, wo", helloRoute, worldRoute, ", wo")
	testRoute_findSubRoute(t, root, "hello, world, again", worldRoute, nil, ", again")
}

func testRoute_findSubRoute(t *testing.T, root *Route, path string, expectedParent, expectedFound *Route, expectedRemainingPath string) {
	parent, found, remainingPath := root.findSubRoute(path)
	if parent != expectedParent || found != expectedFound || remainingPath != expectedRemainingPath {
		t.Errorf("%v.findSubRoute(%q) = %v, %v, %q want %v, %v, %q",
			root, path, parent, found, remainingPath, expectedParent, expectedFound, expectedRemainingPath)
	}
}

//type insertTest struct {
//	path        string
//	resultPath  string
//	resultPaths []string
//}
//
//func TestRoute_insert(t *testing.T) {
//	root := newRoute("root")
//	emptyResult := root.insert("")
//	if emptyResult != root {
//		t.Fail()
//	}
//
//	//from https://en.wikipedia.org/wiki/Radix_tree
//	tests := []insertTest{
//		{"test", "test", []string{"", "test"}},
//		{"slow", "slow", []string{"", "slow", "test"}},
//		{"water", "water", []string{"", "slow", "test", "water"}},
//		{"slower", "er", []string{"", "slow", "slower", "test", "water"}},
//	}
//	testRoute_insert(t, tests)
//
//	tests = []insertTest{
//		{"tester", "tester", []string{"", "tester"}},
//		{"test", "test", []string{"", "test", "tester"}},
//		{"team", "am", []string{"", "te", "team", "test", "tester"}},
//		{"toast", "oast", []string{"", "t", "te", "team", "test", "tester", "toast"}},
//	}
//	testRoute_insert(t, tests)
//
//	tests = []insertTest{
//		{"romane", "romane", []string{"", "romane"}},
//		{"romanus", "us", []string{"", "roman", "romane", "romanus"}},
//		{"romulus", "ulus", []string{"", "rom", "roman", "romane", "romanus", "romulus"}},
//		{"rubens", "ubens", []string{"", "r", "rom", "roman", "romane", "romanus", "romulus", "rubens"}},
//		{"ruber", "r", []string{"", "r", "rom", "roman", "romane", "romanus", "romulus", "rube", "rubens", "ruber"}},
//		{"rubicon", "icon", []string{"", "r", "rom", "roman", "romane", "romanus", "romulus", "rub", "rube", "rubens", "ruber", "rubicon"}},
//		{"rubicundus", "undus", []string{"", "r", "rom", "roman", "romane", "romanus", "romulus", "rub", "rube", "rubens", "ruber", "rubic", "rubicon", "rubicundus"}},
//	}
//	testRoute_insert(t, tests)
//
//	//reverse order of previous suite. any permutation should result in the same trie at the final step.
//	tests = []insertTest{
//		{"rubicundus", "rubicundus", []string{"", "rubicundus"}},
//		{"rubicon", "on", []string{"", "rubic", "rubicon", "rubicundus"}},
//		{"ruber", "er", []string{"", "rub", "ruber", "rubic", "rubicon", "rubicundus"}},
//		{"rubens", "ns", []string{"", "rub", "rube", "rubens", "ruber", "rubic", "rubicon", "rubicundus"}},
//		{"romulus", "omulus", []string{"", "r", "romulus", "rub", "rube", "rubens", "ruber", "rubic", "rubicon", "rubicundus"}},
//		{"romanus", "anus", []string{"", "r", "rom", "romanus", "romulus", "rub", "rube", "rubens", "ruber", "rubic", "rubicon", "rubicundus"}},
//		{"romane", "e", []string{"", "r", "rom", "roman", "romane", "romanus", "romulus", "rub", "rube", "rubens", "ruber", "rubic", "rubicon", "rubicundus"}},
//	}
//	testRoute_insert(t, tests)
//}
//
//func testRoute_insert(t *testing.T, tests []insertTest) {
//	root := newRoute("")
//	for _, test := range tests {
//		result := root.insert(test.path)
//		resultPaths := root.listAllPaths()
//		if result.path != test.resultPath || !reflect.DeepEqual(resultPaths, test.resultPaths) {
//			t.Errorf("insert(%q) = %q, %v want %q, %v", test.path, result.path, resultPaths, test.resultPath, test.resultPaths)
//		}
//	}
//}
//
//func TestRoute_insertChildPath(t *testing.T) {
//	root := newRoute("root")
//	emptyResult := root.SubRoute("")
//	if emptyResult != root {
//		t.Fail()
//	}
//	tests := []struct {
//		path        string
//		resultPath  string
//		resultPaths []string
//	}{
//		{"math", "math", []string{"root", "rootmath"}},
//		{"mathematics", "ematics", []string{"root", "rootmath", "rootmathematics"}},
//		{"maybe", "ybe", []string{"root", "rootma", "rootmath", "rootmathematics", "rootmaybe"}},
//		{"div", "div", []string{"root", "rootdiv", "rootma", "rootmath", "rootmathematics", "rootmaybe"}},
//	}
//	for _, test := range tests {
//		result := root.insertChildPath(test.path)
//		resultPaths := root.listAllPaths()
//		if result.path != test.resultPath || !reflect.DeepEqual(resultPaths, test.resultPaths) {
//			t.Errorf("%v.insertChildPath(%q) = %v, %v want %v %v", root, test.path, result, resultPaths, test.resultPath, test.resultPaths)
//		}
//	}
//}
//
//func TestRoute_splitPathToPrefix_building(t *testing.T) {
//	tests := []struct {
//		path        string
//		prefix      string
//		resultPaths []string
//	}{
//		{"root", "", []string{"root"}},
//		{"root", "ro", []string{"ro", "root"}},
//		{"root", "root", []string{"root"}},
//	}
//	for _, test := range tests {
//		route := newRoute(test.path)
//		route.splitPathToPrefix(test.prefix)
//		resultPaths := route.listAllPaths()
//		if !reflect.DeepEqual(resultPaths, test.resultPaths) {
//			t.Errorf("%v.splitPathToPrefix(%q) = %v want %v", route, test.prefix, resultPaths, test.resultPaths)
//		}
//	}
//}
//
//func TestRoute_splitPathToPrefix_newChild(t *testing.T) {
//	rh := &routeHandler{}
//	teamChildren := []*Route{newRoute("mates"), newRoute("_sub")}
//	team := &Route{"team", teamChildren, rh}
//	team.splitPathToPrefix("te")
//	resultPaths := team.listAllPaths()
//	if !reflect.DeepEqual(resultPaths, []string{"te", "team", "teammates", "team_sub"}) {
//		t.Error("resultPaths are not what they should be")
//	}
//	if len(team.children) != 1 {
//		t.Error("len(team.children) should be 1")
//	}
//	if !areRoutesEqual(team.children[0].children, teamChildren) {
//		t.Error("team's child's children are not what they should be")
//	}
//	if team.children[0].routeHandler != rh {
//		t.Error("teams's child's routeHandler is not what it should be")
//	}
//}
//
//func TestRoute_listAllPaths(t *testing.T) {
//	root := &Route{
//		"a",
//		[]*Route{
//			newRoute("b"),
//			&Route{
//				"c",
//				[]*Route{newRoute("d"), newRoute("e")},
//				nil,
//			},
//		},
//		nil,
//	}
//	tests := []struct {
//		root   *Route
//		result []string
//	}{
//		{root, []string{"a", "ab", "ac", "acd", "ace"}},
//		{newRoute("nil children"), []string{"nil children"}},
//		{&Route{"simple", []*Route{newRoute("a"), newRoute("b")}, nil}, []string{"simple", "simplea", "simpleb"}},
//		{&Route{"a", []*Route{&Route{"b", []*Route{newRoute("c")}, nil}}, nil}, []string{"a", "ab", "abc"}},
//	}
//	for _, test := range tests {
//		result := test.root.listAllPaths()
//		if !reflect.DeepEqual(result, test.result) {
//			t.Errorf("%v.listAllPaths() = %v want %v", test.root, result, test.result)
//		}
//	}
//}
//
//func (route *Route) listAllPaths() []string {
//	result := []string{route.path}
//	for _, child := range route.children {
//		childPaths := child.listAllPaths()
//		for _, childPath := range childPaths {
//			result = append(result, route.path+childPath)
//		}
//	}
//	return result
//}
//
//func TestLevelOrder(t *testing.T) {
//	leafNil := newRoute("leafNil")
//	leafEmpty := &Route{"leafEmpty", []*Route{}, nil}
//	singleChild := &Route{"singleChild", []*Route{division}, nil}
//	multiChild := &Route{"multiChild", []*Route{another, boy, chooses}, nil}
//	bigsChild := &Route{"bigsChild", []*Route{leafEmpty, singleChild}, nil}
//	big := &Route{"big", []*Route{leafNil, bigsChild, multiChild, frogs}, nil}
//	tests := []struct {
//		root   *Route
//		result []*Route
//	}{
//		{leafNil, []*Route{leafNil}},
//		{leafEmpty, []*Route{leafEmpty}},
//		{singleChild, []*Route{singleChild, division}},
//		{bigsChild, []*Route{bigsChild, leafEmpty, singleChild, division}},
//		{multiChild, []*Route{multiChild, another, boy, chooses}},
//		{big, []*Route{big, leafNil, bigsChild, multiChild, frogs, leafEmpty, singleChild, another, boy, chooses, division}},
//	}
//	for _, test := range tests {
//		result := test.root.levelOrder()
//		if !areRoutesEqual(result, test.result) {
//			t.Errorf("%v.levelOrder() = %v want %v", test.root, result, test.result)
//		}
//	}
//}
//
//func (route *Route) levelOrder() []*Route {
//	result := []*Route{}
//	queue := []*Route{route}
//	for len(queue) > 0 {
//		temp := queue[0]
//		result = append(result, temp)
//		queue = append(queue, temp.children...)
//		queue = queue[1:]
//	}
//	return result
//}
//
//func TestRoute_findOrCreateChildWithCommonPrefix(t *testing.T) {
//	tests := []struct {
//		path     string
//		children []*Route
//		child    *Route
//		prefix   string
//		created  bool
//	}{
//		{"", []*Route{}, nil, "", true},
//		{"hello", []*Route{}, nil, "hello", true},
//		{"", []*Route{a, b}, nil, "", true},
//		{"", []*Route{empty, hello}, empty, "", false},
//		{"character", []*Route{another, boy, chooses, division}, chooses, "ch", false},
//		{"divisor", []*Route{another, boy, chooses, division}, division, "divis", false},
//		{"ant", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, another, "an", false},
//		{"ant", []*Route{boy, chooses, division, elephant, frogs, giraffe}, nil, "ant", true},
//		{"hello", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, nil, "hello", true},
//		{"boy", []*Route{another, chooses, division, elephant}, nil, "boy", true},
//		{"boy", []*Route{another, boy, chooses, division, elephant, frogs}, boy, "boy", false},
//		{"boys", []*Route{another, boy, chooses}, boy, "boy", false},
//	}
//	for _, test := range tests {
//		route := &Route{"route", test.children, nil}
//		child, prefix := route.findOrCreateChildWithCommonPrefix(test.path)
//		passed := prefix == test.prefix && child != nil && (test.created || child == test.child)
//		if !passed {
//			t.Errorf("route.findOrCreateChildWithCommonPrefix(%q) = %v, %q want %v, %q (%v)",
//				test.path, child, prefix, test.child, test.prefix, test.created)
//		}
//		if test.created {
//			index, _ := route.indexOfCommonPrefixChild(test.path)
//			if route.children[index] != child {
//				t.Errorf("route.findOrCreateChildWithCommonPrefix() did not appropriately create child")
//			}
//		}
//	}
//}

func TestRoute_findChildWithCommonPrefix(t *testing.T) {
	tests := []struct {
		path     string
		children []*Route
		child    *Route
		index    int
		prefix   string
	}{
		{"", []*Route{}, nil, -1, ""},
		{"hello", nil, nil, -1, ""},
		{"hello", []*Route{}, nil, -1, ""},
		{"", []*Route{a, b}, nil, -1, ""},
		{"", []*Route{empty, hello}, empty, 0, ""},
		{"character", []*Route{another, boy, chooses, division}, chooses, 2, "ch"},
		{"divisor", []*Route{another, boy, chooses, division}, division, 3, "divis"},
		{"ant", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, another, 0, "an"},
		{"ant", []*Route{boy, chooses, division, elephant, frogs, giraffe}, nil, -1, ""},
		{"hello", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, nil, -8, ""},
		{"boy", []*Route{another, chooses, division, elephant}, nil, -2, ""},
		{"boy", []*Route{another, boy, chooses, division, elephant, frogs}, boy, 1, "boy"},
		{"boys", []*Route{another, boy, chooses}, boy, 1, "boy"},
	}
	for _, test := range tests {
		route := newRoute("", test.children...)
		child, index, prefix := route.findChildWithCommonPrefix(test.path)
		if child != test.child || index != test.index || prefix != test.prefix {
			t.Errorf("route.findChildWithCommonPrefix(%q) = %v, %v, %q want %v, %v, %q",
				test.path, child, index, prefix, test.child, test.index, test.prefix)
		}
	}
}

func TestRoute_indexOfCommonPrefixChild(t *testing.T) {
	tests := []struct {
		path     string
		children []*Route
		index    int
		prefix   string
	}{
		{"", []*Route{}, -1, ""},
		{"hello", []*Route{}, -1, ""},
		{"", []*Route{a, b}, -1, ""},
		{"", []*Route{empty, hello}, 0, ""},
		{"character", []*Route{another, boy, chooses, division}, 2, "ch"},
		{"divisor", []*Route{another, boy, chooses, division}, 3, "divis"},
		{"ant", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, 0, "an"},
		{"ant", []*Route{boy, chooses, division, elephant, frogs, giraffe}, -1, ""},
		{"hello", []*Route{another, boy, chooses, division, elephant, frogs, giraffe}, -8, ""},
		{"boy", []*Route{another, chooses, division, elephant}, -2, ""},
		{"boy", []*Route{another, boy, chooses, division, elephant, frogs}, 1, "boy"},
		{"boys", []*Route{another, boy, chooses}, 1, "boy"},
	}
	for _, test := range tests {
		route := newRoute("", test.children...)
		index, prefix := route.indexOfCommonPrefixChild(test.path)
		if index != test.index || prefix != test.prefix {
			t.Errorf("route.indexOfCommonPrefixChild(%q) = %v, %q want %v, %q",
				test.path, index, prefix, test.index, test.prefix,
			)
		}
	}
}

func TestRoute_insertChildAtIndex(t *testing.T) {
	one := newRoute("one")
	two := newRoute("two")
	three := newRoute("three")
	tests := []struct {
		children []*Route
		insert   *Route
		index    int
		result   []*Route
	}{
		{nil, one, -1, nil},
		{nil, one, 2, nil},
		{[]*Route{}, one, 2, []*Route{}},
		{[]*Route{one}, two, 4, []*Route{one}},
		{nil, one, 0, []*Route{one}},
		{[]*Route{}, one, 0, []*Route{one}},
		{[]*Route{one}, two, 0, []*Route{two, one}},
		{[]*Route{one}, two, 1, []*Route{one, two}},
		{[]*Route{one, two}, three, 0, []*Route{three, one, two}},
		{[]*Route{one, two}, three, 1, []*Route{one, three, two}},
		{[]*Route{one, two}, three, 2, []*Route{one, two, three}},
		//these should not occur during normal use, but still testing.
		{[]*Route{one, two}, nil, 1, []*Route{one, nil, two}},
		{[]*Route{one, two}, two, 1, []*Route{one, two, two}},
	}
	for _, test := range tests {
		route := &Route{"route", test.children, nil}
		route.insertChildAtIndex(test.insert, test.index)
		equals := areRoutesEqual(route.children, test.result)
		if !equals {
			t.Errorf("%v insertChildAtIndex(%v, %v) = %v want %v", test.children, test.insert, test.index, route.children, test.result)
		}
	}
}

func TestAreRoutesEqual(t *testing.T) {
	one, two := newRoute("one"), newRoute("two")
	tests := []struct {
		a      []*Route
		b      []*Route
		equals bool
	}{
		{nil, nil, true},
		{nil, []*Route{}, false},
		{[]*Route{}, []*Route{}, true},
		{nil, []*Route{one}, false},
		{[]*Route{one}, []*Route{two}, false},
		{[]*Route{one, two}, []*Route{one, two}, true},
		{[]*Route{one, nil}, []*Route{one, nil}, true},
		{[]*Route{one, nil}, []*Route{one, two}, false},
	}
	for _, test := range tests {
		equals := areRoutesEqual(test.a, test.b)
		if equals != test.equals {
			t.Errorf("areRoutesEqual(%v, %v) = %v want %v", test.a, test.b, equals, test.equals)
		}
	}
}

func areRoutesEqual(a, b []*Route) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
