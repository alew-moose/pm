# Тестовое задание - пакетный менеджер

[Задание](task.md)


- [.] TODO:
  - [ ] readme:
    - [ ] моя интерпретация packets в файле
    - [ ] Makefile
    - [ ] config in HOME
    - [ ] сравнение версий
    - [ ] absolute paths?
    - [ ] нет возможности рекурсивно добавить файлы (double star)
    - [ ] ...
  - [ ] packageDownloader -> downloader && packageUploader -> uploader ?
  - [X] validate config
  - [ ] log verbose:
    - [ ] downloader
    - [ ] uploader
  - [ ] regexp replace all -> strings replace all?
  - [ ] нужны другие варианты подключения по ssh, кроме ssh-agent?
  - [ ] create path unless exists (in constructor?)
  - [ ] проверить, что сохраняются пермишны
  - [ ] print working directory?
  - [ ] move package version & package version spec to pkg? + tests ?
  - [ ] name -> type PackageName + validation + move to pkg?
  - [ ] recursive (double star)
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
  - [X] download packets
  - [X] remove "failed to"
  - [ ] выводить статистику после загрузки/выгрузки? (кол-во файлов, общий размер)
