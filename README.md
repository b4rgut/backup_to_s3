# cloud_backuper

Это утилита для загрузки файлов с определенным префиксом в S3 хранилище.

## Настройка

### Добавление AccessKeyID и SecretAccessKeyID S3 в Windows Credentials:
1. Нажимаем `Win + R`.
2. Запускаем команду `rundll32.exe keymgr.dll,KRShowKeyMgr`.
3. Откроется классический диалоговый интерфейс **Сохранение имен пользователей и паролей (Stored User Names and Password)**. Нажимаем кнопку **Добавить**:

   ![Screenshot from 2024-07-21 02-59-18](https://github.com/user-attachments/assets/0a383905-0592-4963-8551-525bdd846c34)
5. Заполняем поля как на картинке:

   ![Screenshot from 2024-07-21 03-02-09](https://github.com/user-attachments/assets/fd81beff-c4cf-47d8-8ae6-a78cf3cd4ce0)

### Файл конфигурации:

Допустим,  нам нужно выгрузить из папки `D:/files/` файлы с именами формата:
1. `prefix_1_yyyy_mm_dd.txt`
2. `prefix_2_yyyy_mm_dd.txt`
3. `prefix_3_yyyy_mm_dd.txt`

и загрузить их каждый в свою папку в S3 хранилище.

Тогда файл конфигурации будет иметь следующий вид:
```yml
s3:
   endpoint: "my_s3_endpoint.com" 
   backet: "my_s3_backet_name"
 
 local_directory_path: "D:/files/"
 windows_credential: "my_name_windows_credential"
 
 directory_struct:
   - prefix_file: "prefix_1"
     cloud_dir: "prefix_1"
   - prefix_file: "prefix_2"
     cloud_dir: "prefix_2"
   - prefix_file: "prefix_3"
     cloud_dir: "prefix_3"
```

В S3 хранилище мы получим следующую структуру файлов:
```
my_s3_backet_name:
  - prefix-1:
    - prefix_1_yyyy_mm_dd.txt
  - prefix-2:
    - prefix_2_yyyy_mm_dd.txt`
  - prefix-3:
    - prefix_1_yyyy_mm_dd.txt`
```

**Примечание:** 
1. Файлы, которые не имеют нужного префикса, игнорируются при загрузке!
2. Можно указывать произвольное количество префиксов файлов.

## Сборка

```
make
```



   

