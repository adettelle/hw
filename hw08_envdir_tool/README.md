## Домашнее задание №8 «Утилита envdir»

Необходимо реализовать утилиту `envdir` на Go.

С помощью утилиты запускаются программы, для которых формируются переменные окружения из файлов расположенных в каталоге.

Пример вызова:

```bash
$ go-envdir /path/to/env/dir command arg1 arg2
```
*go-envdir* - утилита  
*/path/to/env/dir* - каталог с файлами переменных окружения  
*command* - запускаемая утилитой программа  
*arg1 arg2* - параметры запускаемой программы  

### Требования:
- утилита `envdir` читает первые строки `T`, файлов с именем `S` из каталога, переданного утилите в параметрах
- `envdir` удаляет переменную среды с именем `S`, если таковая существует
- `envdir` добавляет переменную среды с именем `S` и значением `T`, если файл не пустой;
- имя `S` не должно содержать `=`;
- пробелы и табуляция в конце `T` удаляются;
- терминальные нули в `T` (`0x00`) заменяются на перевод строки (`\n`);
- стандартные потоки ввода/вывода/ошибок пробрасывались в вызываемую программу;
- код выхода утилиты должен совпадать с кодом выхода программы.

При необходимости можно выделять дополнительные функции / ошибки.

Юнит-тесты могут использовать файлы из `testdata` или создавать свои директории / файлы,
которые **обязаны** подчищать после своего выполнения.

---
Пример использования:
```bash
$ go-envdir /path/to/env/dir command arg1 arg2
```
Если в директории `/path/to/env/dir` содержатся файлы:
* `FOO` с содержимым `123`;
* `BAR` с содержимым `value`,

то вызов выше эквивалентен вызову
```bash
$ FOO=123 BAR=value command arg1 arg2
```
---


### Критерии оценки
- Пайплайн зелёный - 4 балла
- Добавлены юнит-тесты - до 4 баллов
- Понятность и чистота кода - до 2 баллов

#### Зачёт от 7 баллов

### Подсказки
- https://www.unix.com/man-page/debian/8/envdir/
- `os.Args`
- `os.ReadDir`
- `bytes.Replace`, `strings.TrimRight`
- `exec.Command`
