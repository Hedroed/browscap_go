package browscap_go

type dictionary struct {
	mapped map[string]section

	completeData map[string]map[string]string

	tree *ExpressionTree
}

type section map[string]string

func newDictionary() *dictionary {
	return &dictionary{
		mapped: make(map[string]section),

		completeData: make(map[string]map[string]string),

		tree: NewExpressionTree(),
	}
}

func (dict *dictionary) buildCompleteData() {
	for name, _ := range dict.mapped {
		dict.completeData[name] = dict.buildData(name)
	}
}

func (dict *dictionary) getData(name string) map[string]string {
	return dict.completeData[name]
}

func (dict *dictionary) buildData(name string) map[string]string {
	res := make(map[string]string)

	if item, found := dict.mapped[name]; found {
		// Parent's data
		if parentName, hasParent := item["Parent"]; hasParent {
			parentData := dict.buildData(parentName)
			if len(parentData) > 0 {
				for k, v := range parentData {
					if k == "Parent" {
						continue
					}
					res[k] = v
				}
			}
		}
		// It's item data
		if len(item) > 0 {
			for k, v := range item {
				if k == "Parent" {
					continue
				}
				res[k] = v
			}
		}
	}

	return res
}
