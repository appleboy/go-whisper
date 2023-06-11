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

For more information see [ggerganov/whisper.cpp][3].

[2]: https://huggingface.co/ggerganov/whisper.cpp/tree/main
[3]: https://github.com/ggerganov/whisper.cpp/tree/master/models

## Prepare

Download the model you want to use and put it in the `models` directory.

```sh
curl -LJ https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin \
  --output models/ggml-small.bin
```

## Usage

Please follow these simplified instructions to transcribe the audio file using a Docker container:

1. Ensure that you have a `testdata` directory containing the `jfk.wav` file.
2. Mount both the `models` and `testdata` directories to the Docker container.
3. Specify the model using the `--model` flag and the audio file path using the `--audio-path` flag.
4. The transcript result file will be saved in the same directory as the audio file.

To transcribe the audio file, execute the command provided below.

```sh
docker run \
  -v $PWD/models:/app/models \
  -v $PWD/testdata:/app/testdata \
  ghcr.io/appleboy/go-whisper:latest \
  --model /app/models/ggml-small.bin \
  --audio-path /app/testdata/jfk.wav
```

See the following output:

```sh
whisper_init_from_file_no_state: loading model from '/app/models/ggml-small.bin'
whisper_model_load: loading model
whisper_model_load: n_vocab       = 51865
whisper_model_load: n_audio_ctx   = 1500
whisper_model_load: n_audio_state = 768
whisper_model_load: n_audio_head  = 12
whisper_model_load: n_audio_layer = 12
whisper_model_load: n_text_ctx    = 448
whisper_model_load: n_text_state  = 768
whisper_model_load: n_text_head   = 12
whisper_model_load: n_text_layer  = 12
whisper_model_load: n_mels        = 80
whisper_model_load: ftype         = 1
whisper_model_load: qntvr         = 0
whisper_model_load: type          = 3
whisper_model_load: mem required  =  743.00 MB (+   16.00 MB per decoder)
whisper_model_load: adding 1608 extra tokens
whisper_model_load: model ctx     =  464.68 MB
whisper_model_load: model size    =  464.44 MB
whisper_init_state: kv self size  =   15.75 MB
whisper_init_state: kv cross size =   52.73 MB
1:46AM INF system_info: n_threads = 8 / 8 | AVX = 0 | AVX2 = 0 | AVX512 = 0 | FMA = 0 | NEON = 1 | ARM_FMA = 1 | F16C = 0 | FP16_VA = 0 | WASM_SIMD = 0 | BLAS = 0 | SSE3 = 0 | VSX = 0 | COREML = 0 | 
 module=transcript
whisper_full_with_state: auto-detected language: en (p = 0.967331)
1:46AM INF [    0s ->    11s] And so my fellow Americans, ask not what your country can do for you, ask what you can do for your country. module=transcript
```
