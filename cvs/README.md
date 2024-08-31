# Container Versioning Systems

## Docker Save
I exported the image of the server from the development assignment and analysed it. I found the following files:-
1. index.json :- Contains the image index which includes annotations in the Open Container Initiative (OCI) format.
2. manifest.json :- Contains a list of the layers and blobs in the image, also contains some meta data
                    (like tags or the hash of the config blob).
3. oci-layout:- Contains the version of OCI layout specification.
4. repositories:- Contains the name of the image and its corresponding layer hash
5. blobs/ :- Contains sha256/ which contains the blobs corresponding to different layers (or manifests or config) 
                and the filenames are the sha256 hashes of the blobs.
6. On extracting the tar (blob) corresponding to a layer, we find a manifest of files written to/deleted (white-outs) in that layer.


