package level

import "math/rand/v2"

func shuffleSlice[T any](random *rand.Rand, slice []T) {
	random.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
