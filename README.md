# Тестовое задание - пакетный менеджер

[Задание](task.md)


- [.] TODO:
  - [ ] packageDownloader -> downloader && packageUploader -> uploader ?
  - [X] validate config
  - [ ] regexp replace all -> strings replace all?
  - [ ] create path unless exists (in constructor?)
  - [ ] move package version & package version spec to pkg? + tests ?
  - [ ] name -> type PackageName + validation + move to pkg?
  - [X] что означает "packets" в файле для упаковки packet.json?
  - [ ] tests for uploader/downloader configs from json/yaml
  - [X] dedup files for uploader
  - [ ] pretty printer для слайса стрингеров
  - [ ] разнобой с методами по значению/по указателю
  - [ ] check tar ErrInsecurePath
  - [ ] test absolute paths
  - [ ] check filepath.IsAbs
  - [ ] integration tests
  - [ ] shorten packageVersion -> pv & packageVersionSpec -> pvs ?
  - [ ] sftpClient refactor : remotePath method
  - [ ] сделать архиватор отдельно (упаковка, распаковка)
  - [ ] enable more linters? https://golangci-lint.run/docs/linters/configuration/
  - [ ] там, где var smth, var error: может, сделать именованные возвращаемые?
  - [ ] https://sftptogo.com/blog/go-sftp/ get host key
  - [ ] валидации наслаиваются, кажется по нескольку раз вызываю
  - [ ] %q -> '%s' ?
  - [ ] неэффективно хранить все пакеты в одной директории
  - [ ] uploader:
    - [ ] download packets
  - [ ] remove "failed to"
  - [ ] выводить статистику после загрузки/выгрузки? (кол-во файлов, общий размер)
