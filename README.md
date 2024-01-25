# go-directory-logger

Утилита, следящая за изменениями файлов в указанной папке и во всех внутренних папках. Все изменения записываются в базу данных MySQL

## Инструкция по запуску:
1. Созать БД MySQL с таблицей:
```sql
CREATE TABLE files (
  id INT NOT NULL AUTO_INCREMENT,
  dirPath VARCHAR(255) NOT NULL,
  filename VARCHAR(255) NOT NULL,
  operation VARCHAR(255) NOT NULL,
  date TIMESTAMP,
  PRIMARY KEY(id)
  );
```
2. В файле ```config.yml``` запишите все необходимые данные о базе данных и отслеживаемых папках
3. Установить все необходимые пакеты
4. Запустить ```cmd/main.go```
