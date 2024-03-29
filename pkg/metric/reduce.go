package metric

func Sum(iter BucketIterator) float64 {
	var result = 0.0
	for iter.Next() {
		bucket := iter.Bucket()
		for _, p := range bucket.Points {
			result += p
		}
	}
	return result
}

func Avg(iter BucketIterator) float64 {
	var result = 0.0
	var count = 0.0
	for iter.Next() {
		bucket := iter.Bucket()
		for _, p := range bucket.Points {
			result += p
			count++
		}
	}
	return result / count
}

func Min(iterator BucketIterator) float64 {
	var result = 0.0
	var started = false
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !started {
				result = p
				started = true
				continue
			}
			if p < result {
				result = p
			}
		}
	}
	return result
}

func Max(iterator BucketIterator) float64 {
	var result = 0.0
	var started = false
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !started {
				result = p
				started = true
				continue
			}
			if p > result {
				result = p
			}
		}
	}
	return result
}

func Count(iterator BucketIterator) float64 {
	var result int64
	for iterator.Next() {
		bucket := iterator.Bucket()
		result += bucket.Count
	}
	return float64(result)
}
