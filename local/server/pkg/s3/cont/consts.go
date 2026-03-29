package cont

const (
	// HashPath defines the storage path for hash data within the 'FDIM' directory.
	hashPath = "FDIM/data/hash/"

	// TempPath specifies the directory for temporary files in the 'FDIM' structure.
	tempPath = "FDIM/temp/"

	// DirectPath indicates the directory for direct uploads or access within the 'FDIM' structure.
	DirectPath = "FDIM/direct"

	// UploadTypeMultipart represents the identifier for multipart uploads,
	// allowing large files to be uploaded in chunks.
	UploadTypeMultipart = 1

	// UploadTypePresigned signifies the use of presigned URLs for uploads,
	// facilitating secure, authorized file transfers without requiring direct access to the storage credentials.
	UploadTypePresigned = 2

	// PartSeparator is used as a delimiter in multipart upload processes,
	// separating individual file parts.
	partSeparator = ","
)
