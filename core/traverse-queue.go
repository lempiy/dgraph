package core

type TraverseQueue struct {
	s []*NodeOutput
}

func NewTraverseQueue() *TraverseQueue {
	return new(TraverseQueue)
}

func (q *TraverseQueue) add(incomeId *string, bufferQueue *TraverseQueue, items ...*NodeInput) {
	for _ ,itm := range items {
		item := q.find(func(v *NodeOutput) bool {
			return v.Id == itm.Id
		})
		if item == nil && bufferQueue != nil {
			item = bufferQueue.find(func(v *NodeOutput) bool {
				return v.Id == itm.Id
			})
		}
		if item != nil && incomeId != nil {
			item.PassedIncomes = append(item.PassedIncomes, *incomeId)
			return
		}
		var incomes []string
		if incomeId != nil {
			incomes = append(incomes, *incomeId)
		}
		q.s = append(q.s, &NodeOutput{
			NodeInput: &NodeInput{
				Id: itm.Id,
				Next: itm.Next,
			},
			PassedIncomes: append([]string{}, incomes...),
			RenderIncomes: append([]string{}, incomes...),
			ChildrenOnMatrix: 0,
		})
	}
}

func (q *TraverseQueue) find(cb func(v *NodeOutput)bool)*NodeOutput {
	for _, itm := range q.s {
		if cb(itm) {
			return itm
		}
	}
	return nil
}

func (q *TraverseQueue) push(item *NodeOutput) {
	q.s = append(q.s, item)
}

func (q *TraverseQueue) length() int {
	return len(q.s)
}

func (q *TraverseQueue) some(cb func(v *NodeOutput)bool) bool {
	for _, itm := range q.s {
		if cb(itm) {
			return true
		}
	}
	return false
}

func (q *TraverseQueue) shift() *NodeOutput {
	if len(q.s) == 0 {
		return nil
	}
	item := *q.s[0]
	q.s = append(q.s, q.s[1:]...)
	return &item
}

func (q *TraverseQueue) drain() *TraverseQueue {
	newQ := NewTraverseQueue()
	newQ.s = append(newQ.s, q.s...)
	q.s = []*NodeOutput{}
	return newQ
}
