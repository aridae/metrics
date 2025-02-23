package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
)

// standardPasses возвращает все стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
func standardPasses() []*analysis.Analyzer {
	return _allStandardPassesAnalyzers
}

// _allStandardPassesAnalyzers все стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
var _allStandardPassesAnalyzers = []*analysis.Analyzer{
	appends.Analyzer,             // Проверка корректности использования append() для срезов
	asmdecl.Analyzer,             // Проверка объявлений ассемблерного кода
	assign.Analyzer,              // Проверка присваиваний значений переменным
	atomic.Analyzer,              // Проверка использования атомарных операций
	atomicalign.Analyzer,         // Проверка выравнивания данных для атомарных операций
	bools.Analyzer,               // Проверка использования булевых выражений
	buildssa.Analyzer,            // Построение SSA-формы программы
	buildtag.Analyzer,            // Проверка использования build tags
	cgocall.Analyzer,             // Проверка вызовов C-функций
	composite.Analyzer,           // Проверка составных литералов
	copylock.Analyzer,            // Проверка копирования заблокированных ресурсов
	ctrlflow.Analyzer,            // Проверка управления потоками выполнения
	deepequalerrors.Analyzer,     // Проверка ошибок сравнения глубоких структур
	defers.Analyzer,              // Проверка использования defer
	directive.Analyzer,           // Проверка директив компилятора
	errorsas.Analyzer,            // Проверка преобразования ошибок
	fieldalignment.Analyzer,      // Проверка выравнивания полей структуры
	findcall.Analyzer,            // Поиск вызовов функций
	framepointer.Analyzer,        // Проверка использования указателя кадра стека
	httpmux.Analyzer,             // Проверка маршрутизации HTTP-запросов
	httpresponse.Analyzer,        // Проверка обработки HTTP-ответов
	ifaceassert.Analyzer,         // Проверка утверждений интерфейсов
	inspect.Analyzer,             // Инспекция кода
	loopclosure.Analyzer,         // Проверка замыканий в циклах
	lostcancel.Analyzer,          // Проверка потери контекста отмены
	nilfunc.Analyzer,             // Проверка использования nil-функций
	nilness.Analyzer,             // Проверка использования nil-значений
	pkgfact.Analyzer,             // Проверка фактов пакетов
	printf.Analyzer,              // Проверка форматированной печати
	reflectvaluecompare.Analyzer, // Проверка сравнений отражаемых значений
	shadow.Analyzer,              // Проверка затенения имен переменных
	shift.Analyzer,               // Проверка сдвигов битовых значений
	sigchanyzer.Analyzer,         // Проверка использования каналов сигналов
	slog.Analyzer,                // Проверка использования журналирования
	sortslice.Analyzer,           // Проверка сортировки срезов
	stdmethods.Analyzer,          // Проверка методов стандартного пакета
	stdversion.Analyzer,          // Проверка совместимости с версией стандарта
	stringintconv.Analyzer,       // Проверка преобразований строк в целые числа
	structtag.Analyzer,           // Проверка тегов структур
	testinggoroutine.Analyzer,    // Проверка горутин в тестах
	tests.Analyzer,               // Проверка тестов
	timeformat.Analyzer,          // Проверка форматирования времени
	unmarshal.Analyzer,           // Проверка распаковки данных
	unreachable.Analyzer,         // Проверка недостижимых участков кода
	unsafeptr.Analyzer,           // Проверка небезопасной работы с указателями
	unusedresult.Analyzer,        // Проверка неиспользуемых результатов функций
	unusedwrite.Analyzer,         // Проверка неиспользуемого записи
	usesgenerics.Analyzer,        // Проверка использования обобщений
	waitgroup.Analyzer,           // Проверка использования wait groups
}
