// Package rag provides retrieval-augmented generation capabilities.
package rag

// MergeResults combines multiple search result slices into a single slice.
// Results are appended in order; no sorting or deduplication is performed.
// Use DeduplicateResults after merging if deduplication is needed.
//
// Example:
//
//	denseResults, _ := vectorDB.Search(ctx, embeddings, 10)
//	sparseResults := bm25Index.Search(ctx, query, 10)
//	merged := rag.MergeResults(denseResults, sparseResults)
//	deduped := rag.DeduplicateResults(merged, func(r SearchResult) string {
//	    return fmt.Sprintf("%d", r.ID)
//	})
func MergeResults(resultSets ...[]SearchResult) []SearchResult {
	var total int
	for _, rs := range resultSets {
		total += len(rs)
	}
	merged := make([]SearchResult, 0, total)
	for _, rs := range resultSets {
		merged = append(merged, rs...)
	}
	return merged
}

// DeduplicateResults removes duplicate results from a merged result set.
// The keyFunc extracts a unique identifier from each result. When multiple
// results share the same key, only the first occurrence (preserving the
// caller's priority ordering) is kept.
//
// This is useful in hybrid search pipelines where the same document may
// appear in both dense and sparse result sets.
//
// Example:
//
//	deduped := rag.DeduplicateResults(merged, func(r SearchResult) string {
//	    if name, ok := r.Fields["name"].(string); ok {
//	        return name
//	    }
//	    return fmt.Sprintf("%d", r.ID)
//	})
func DeduplicateResults(results []SearchResult, keyFunc func(SearchResult) string) []SearchResult {
	seen := make(map[string]bool, len(results))
	deduped := make([]SearchResult, 0, len(results))
	for _, r := range results {
		key := keyFunc(r)
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, r)
		}
	}
	return deduped
}
