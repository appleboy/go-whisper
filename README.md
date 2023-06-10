# go-whisper

docker image for speech to text using [ggerganov/whisper.cpp][1]

[1]:https://github.com/ggerganov/whisper.cpp

## OpenAI's Whisper models converted to ggml format

See the [Available models][2].

| Model      | Disk    | Mem       | SHA                                          |
|------------|---------|-----------|----------------------------------------------|
| tiny       | 75 MB   | ~390 MB   | bd577a113a864445d4c299885e0cb97d4ba92b5f     |
| tiny.en    | 75 MB   | ~390 MB   | c78c86eb1a8faa21b369bcd33207cc90d64ae9df     |
| base       | 142 MB  | ~500 MB   | 465707469ff3a37a2b9b8d8f89f2f99de7299dac     |
| base.en    | 142 MB  | ~500 MB   | 137c40403d78fd54d454da0f9bd998f78703390c     |
| small      | 466 MB  | ~1.0 GB   | 55356645c2b361a969dfd0ef2c5a50d530afd8d5     |
| small.en   | 466 MB  | ~1.0 GB   | db8a495a91d927739e50b3fc1cc4c6b8f6c2d022     |
| medium     | 1.5 GB  | ~2.6 GB   | fd9727b6e1217c2f614f9b698455c4ffd82463b4     |
| medium.en  | 1.5 GB  | ~2.6 GB   | 8c30f0e44ce9560643ebd10bbe50cd20eafd3723     |
| large-v1   | 2.9 GB  | ~4.7 GB   | b1caaf735c4cc1429223d5a74f0f4d0b9b59a299     |
| large      | 2.9 GB  | ~4.7 GB   | 0f4c8e34f21cf1a914c59d8b3ce882345ad349d6     |

For more information see [ggml/whisper][3].

[2]: https://huggingface.co/ggerganov/whisper.cpp/tree/main
[3]: https://github.com/ggerganov/whisper.cpp/tree/master/models
