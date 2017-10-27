package di

import "testing"

func TestLifetime(t *testing.T) {
	getValues := func(l Lifetime, t *testing.T) (int, int, IResolver) {
		defs := NewDefs()
		err := defs.Add(NewA, l)

		if err != nil {
			t.Fatal(err)
		}

		err = defs.Add(NewB, PerDependency)

		if err != nil {
			t.Fatal(err)
		}

		resolver, err := NewResolver(defs)

		if err != nil {
			t.Fatal(err)
		}

		var b B
		resolveErr := resolver.Resolve(&b)

		if resolveErr != nil {
			t.Fatal(resolveErr)
		}

		a1, a2 := b.B()
		return a1, a2, resolver
	}

	t.Run("Singleton", func(t *testing.T) {
		expectedA1 := aCounter + 1
		expectedA2 := expectedA1

		a1, a2, _ := getValues(Singleton, t)

		if a1 != expectedA1 {
			t.Fatal(Singleton, a1, expectedA1)
		}

		if a2 != expectedA2 {
			t.Fatal(Singleton, a2, expectedA2)
		}
	})

	t.Run("PerDependency", func(t *testing.T) {
		expectedA1 := aCounter + 1
		expectedA2 := expectedA1 + 1

		a1, a2, _ := getValues(PerDependency, t)

		if a1 != expectedA1 {
			t.Fatal(Singleton, a1, expectedA1)
		}

		if a2 != expectedA2 {
			t.Fatal(Singleton, a2, expectedA2)
		}
	})

	// per http request handled in another test file

	t.Run("PerResolution", func(t *testing.T) {
		expectedA1 := aCounter + 1
		expectedA2 := expectedA1

		a1, a2, resolver := getValues(PerResolve, t)

		if a1 != expectedA1 {
			t.Fatal(Singleton, a1, expectedA1)
		}

		if a2 != expectedA2 {
			t.Fatal(Singleton, a2, expectedA2)
		}

		expectedA1 = expectedA2 + 1
		expectedA2 = expectedA1

		var b B
		err := resolver.Resolve(&b)

		if err != nil {
			t.Fatal(err)
		}

		a1, a2 = b.B()
		if a1 != expectedA1 {
			t.Fatal(Singleton, a1, expectedA1)
		}

		if a2 != expectedA2 {
			t.Fatal(Singleton, a2, expectedA2)
		}
	})
}
