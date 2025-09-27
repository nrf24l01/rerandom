package task

func (s *Sheet) Rebuild() {
	s.ModifiedUsers = s.ClearUsers
	s.UpdateFractionsAndLive()

	var newRows []SheetRow

	prev_fraction := uint(0)
	for _, row := range s.ModifiedUsers {
		if !row.Alive {
			row.Fraction = 0
		}
		newRows = append(newRows, SheetRow{
			Id:           row.Id,
			Fraction:     row.Fraction,
			FractionFrom: prev_fraction + 1,
			FractionTo:   prev_fraction + row.Fraction,
			LastName:     row.LastName,
			FirstName:    row.FirstName,
			Alive:        row.Alive,
		})
		if row.Alive {
			prev_fraction += row.Fraction
		}
	}
	s.ModifiedUsers = newRows
}

func (s *Sheet) UpdateFractionsAndLive() {
	for _, action := range s.Actions {
		for i, row := range s.ModifiedUsers {
			if row.Id == action.RowId {
				if action.Type == 1 {
					s.ModifiedUsers[i].Fraction = action.Param
				} else if action.Type == 2 {
					if action.Param == 0 {
						s.ModifiedUsers[i].Alive = false
					} else {
						s.ModifiedUsers[i].Alive = true
					}
				}
				break
			}
		}
	}
}

func (s *Sheet) GetTotalFraction() uint {
	var total uint = 0
	for _, row := range s.ModifiedUsers {
		total += row.Fraction
	}
	return total
}

func (s *Sheet) ChangeFraction(userId uint, newFraction uint) {
	s.Actions = append(s.Actions, Action{Type: 1, RowId: userId, Param: newFraction})
	s.Rebuild()
}

func (s *Sheet) ChangeAlive(userId uint, alive bool) {
	var aliveParam uint = 0
	if alive {
		aliveParam = 1
	}
	s.Actions = append(s.Actions, Action{Type: 2, RowId: userId, Param: aliveParam})
	s.Rebuild()
}