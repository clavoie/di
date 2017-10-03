package di

import "testing"

func TestDefs(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		t.Run("Invalid", func(t *testing.T) {
			t.Run("NotFn", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add("invalid", Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("NoOut", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() {}, Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("TooManyOut", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() (int, int, int) { return 1, 2, 3 }, Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("Out1IsErr", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() error { return nil }, Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("Out1IsNotIface", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() int { return 0 }, Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("Out2IsNotErr", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() (B, int) { return nil, 4 }, Singleton)

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("Lifetime", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(func() B { return nil }, (Lifetime)(-1))

				if err == nil {
					t.Fatal(err)
				}
			})
			t.Run("DuplicateDefinition", func(t *testing.T) {
				defs := NewDefs()
				err := defs.Add(NewA, PerResolve)

				if err != nil {
					t.Fatal(err)
				}

				err = defs.Add(NewA, PerResolve)
				if err == nil {
					t.Fatal(err)
				}
			})
		})
	})

	t.Run("Join", func(t *testing.T) {
		defs1 := NewDefs()
		err := defs1.Add(NewA, Singleton)

		if err != nil {
			t.Fatal(err)
		}

		defs2 := NewDefs()
		err = defs2.Add(NewB, Singleton)

		if err != nil {
			t.Fatal(err)
		}

		defs3 := defs1.Join(defs2)
		err = defs3.Add(NewE, Singleton)

		if err != nil {
			t.Fatal(err)
		}

		allDefs := defs3.all()
		hasA := false
		hasB := false
		hasE := false

		if len(allDefs) != 3 {
			t.Fatal(allDefs)
		}

		for _, node := range allDefs {
			switch node.Type {
			case aType:
				hasA = true
			case bType:
				hasB = true
			case eType:
				hasE = true
			}
		}

		if hasA && hasB && hasE {
			return
		}

		t.Fatal(hasA, hasB, hasE)
	})

	t.Run("build", func(t *testing.T) {
		t.Run("cycle", func(t *testing.T) {
			defs := NewDefs()
			err := defs.Add(NewC, Singleton)

			if err != nil {
				t.Fatal(err)
			}

			err = defs.Add(NewD, Singleton)
			if err != nil {
				t.Fatal(err)
			}

			err = defs.Add(NewE, Singleton)
			if err != nil {
				t.Fatal(err)
			}

			_, err = defs.build()
			if err == nil {
				t.Fatal("expecting circular reference err")
			}
		})
		t.Run("DuplicateDefs", func(t *testing.T) {
			defs1 := NewDefs()
			err := defs1.Add(NewA, Singleton)

			if err != nil {
				t.Fatal(err)
			}

			defs2 := NewDefs()
			err = defs2.Add(NewA, PerResolve)

			if err != nil {
				t.Fatal(err)
			}

			defs3 := defs1.Join(defs2)
			_, err = defs3.build()

			if err == nil {
				t.Fatal("no deuplicate definition err")
			}
		})
	})
}
