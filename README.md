# siot-mgo-pool

A mongodb pool by golang

## uri example

mongodb://yapi:abcd1234@vm:27017/yapi
mongodb://yapi:abcd1234@vm:27017,yapi:abcd1234@vm:27017,yapi:abcd1234@vm:27017/yapi

## index options 

### 字段
```golang
// If true, the index will be built in the background on the server and will not block other tasks. The default
// value is false.
Background *bool

// The length of time, in seconds, for documents to remain in the collection. The default value is 0, which means
// that documents will remain in the collection until they're explicitly deleted or the collection is dropped.
ExpireAfterSeconds *int32

// The name of the index. The default value is "[field1]_[direction1]_[field2]_[direction2]...". For example, an
// index with the specification {name: 1, age: -1} will be named "name_1_age_-1".
Name *string

// If true, the index will only reference documents that contain the fields specified in the index. The default is
// false.
Sparse *bool

// Specifies the storage engine to use for the index. The value must be a document in the form
// {<storage engine name>: <options>}. The default value is nil, which means that the default storage engine
// will be used. This option is only applicable for MongoDB versions >= 3.0 and is ignored for previous server
// versions.
StorageEngine interface{}

// If true, the collection will not accept insertion or update of documents where the index key value matches an
// existing value in the index. The default is false.
Unique *bool

// The index version number, either 0 or 1.
Version *int32

// The language that determines the list of stop words and the rules for the stemmer and tokenizer. This option
// is only applicable for text indexes and is ignored for other index types. The default value is "english".
DefaultLanguage *string

// The name of the field in the collection's documents that contains the override language for the document. This
// option is only applicable for text indexes and is ignored for other index types. The default value is the value
// of the DefaultLanguage option.
LanguageOverride *string

// The index version number for a text index. See https://docs.mongodb.com/manual/core/index-text/#text-versions for
// information about different version numbers.
TextVersion *int32

// A document that contains field and weight pairs. The weight is an integer ranging from 1 to 99,999, inclusive,
// indicating the significance of the field relative to the other indexed fields in terms of the score. This option
// is only applicable for text indexes and is ignored for other index types. The default value is nil, which means
// that every field will have a weight of 1.
Weights interface{}

// The index version number for a 2D sphere index. See https://docs.mongodb.com/manual/core/2dsphere/#dsphere-v2 for
// information about different version numbers.
SphereVersion *int32

// The precision of the stored geohash value of the location data. This option only applies to 2D indexes and is
// ignored for other index types. The value must be between 1 and 32, inclusive. The default value is 26.
Bits *int32

// The upper inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
// is ignored for other index types. The default value is 180.0.
Max *float64

// The lower inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
// is ignored for other index types. The default value is -180.0.
Min *float64

// The number of units within which to group location values. Location values that are within BucketSize units of
// each other will be grouped in the same bucket. This option is only applicable to geoHaystack indexes and is
// ignored for other index types. The value must be greater than 0.
BucketSize *int32

// A document that defines which collection documents the index should reference. This option is only valid for
// MongoDB versions >= 3.2 and is ignored for previous server versions.
PartialFilterExpression interface{}

// The collation to use for string comparisons for the index. This option is only valid for MongoDB versions >= 3.4.
// For previous server versions, the driver will return an error if this option is used.
Collation *Collation

// A document that defines the wildcard projection for the index.
WildcardProjection interface{}
```

### 方法
```golang
SetBackground(background bool)
SetExpireAfterSeconds(seconds int32)
SetName(name string)
SetSparse(sparse bool)
SetStorageEngine(engine interface{})
SetUnique(unique bool)
SetVersion(version int32)
SetDefaultLanguage(language string)
SetLanguageOverride(override string)
SetTextVersion(version int32)
SetWeights(weights interface{})
SetSphereVersion(version int32)
SetBits(bits int32)
SetMax(max float64)
SetMin(min float64)
SetBucketSize(bucketSize int32)
SetPartialFilterExpression(expression interface{})
SetCollation(collation *Collation)
SetWildcardProjection(wildcardProjection interface{})
```