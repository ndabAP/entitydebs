package tokenize

// Features represents the features during the tokenization process.
type Features int

const (
	FeatureAll Features = FeatureSyntax | FeatureSentiment

	// FeatureSyntax enables syntax analysis.
	FeatureSyntax Features = 1 << iota
	// FeatureSentiment enables sentiment analysis.
	FeatureSentiment
)
