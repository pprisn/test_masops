package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/template"
)

func main() {
	var err error

	jsonStr := []byte(`
{
"Title": "Заголовок страници",
"H2": ["Строка 1 заголовка","Строка 2 заголовка"],
"BtnInsert": "&#9997;ДОБАВИТЬ НОВУЮ ЗАПИСЬ В ЖУРНАЛ",
"Hcsv": "ID;Time;Name;StatusNSI;Statussdo;Statusupd;Statusauth",
"Tabheader": [
 {"id": "Id"},
 {"updated": "Время обновления"},
 {"name": "Имя"},
 {"status": "Russian Post EAS nsi"},
 {"statussdo": "Russian Post EAS sdo"},
 {"statusupd": "Russian Post EAS Configuration"},
 {"statusauth": "Russian Post EAS user"},
 {"th": ""}
 ],
 "Tabtd": [
 {"id": "ID"},
 {"updated": "UpdatedAt"},
 {"name": "Name"},
 {"status": "Status"},
 {"statussdo": "Statussdo"},
 {"statusupd": "Statusupd"},
 {"statusauth": "Statusauth"}
 ],
 "TitleFormCreate": "Добавление новой записи для мониторинга",
 "TitleFormEit": "Добавление новой записи для мониторинга",
 "Columns": "[0, 1 ,2, 3, 4, 5, 6]",
 "Cntcolumns": 7
 }
`)

	text1 := `
<!DOCTYPE html>
<html>
    <head>
     <meta charset="UTF-8">
     <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">     
     <title>{{$.Title}}</title>
     <!-- Bootstrap CSS -->
     <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
     <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"></script>
     <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>
     <link rel="stylesheet" href="static/bootstrap-4.3.1-dist/css/bootstrap.min.css">
     <script type="text/javascript" src="static/jquery/jquery-3.4.1.min.js"></script>
     <script type="text/javascript" src="static/popper/popper.min.js"></script>
     <script type="text/javascript" src="static/bootstrap-4.3.1-dist/js/bootstrap.min.js"></script>
     <link rel="stylesheet" type="text/css" href="static/css/style.css"> 
     <!-- Full local -->
     <link rel="stylesheet" type="text/css" href="static/DataTables/datatables.min.css"/>
     <script type="text/javascript" src="static/DataTables/datatables.min.js"></script>
    </head>
    <body>
       <div id="overlay"></div><!-- Пoдлoжкa, oднa нa всю стрaницу -->                   
       <h2 class="text-primary" >{{$.H2}}</h2>
       <p> <button type="button" class="btn btn-light" data-toggle="modal" data-target="#insertModal" 
                data-whatever="@Какието данные">{{.BtnInsert}}/button> </p>
       <table id="myTable" class="cell-border compact stripe  buttons"> <!-- responsive -->
       <thead>
    {{ range $_, $v := $.Tabheader }}
		{{- range $_, $v := $v -}}
    {{"\t"}}<th>{{ $v }}</th>{{"\n"}}
		{{- end -}}
    {{ end }}
       </thead>`
	text2 := `	{{range . }}
	<tr data-tr-id={{.ID}}>
`

	text3 := `
    {{ range $_, $v := $.Tabtd }}
		{{- range $k, $v := $v -}}
	   {{"\t "}}  <td class="{{ $k }}">{{"{{"}} .{{ $v }}{{"}}"}}</td>{{"\n"}}
		{{- end -}}
    {{ end }}
    <td><a href="#" class="edit_modal" data-toggle="modal" data-target="#editModal">&#9998;</a></td>`

	text4 := `	</tr>
	{{end}}
	</table>`

	text5 := `<!-- CREATE NEW -->
<div class="modal fade" id="insertModal" tabindex="-1" role="dialog" aria-labelledby="insertModalLabel" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="insertModalLabel">{{.TitleFormCreate}}/h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form id="ajaxCreateForm" method="POST" action='mcreate' >
          <div class="form-group">
   {{ range $_, $v := $.Tabheader }}
		{{- range $k, $v := $v -}}
          {{"\t "}}<label for="recipient-{{$k}}" class="col-form-label">{{ $v }}:</label> {{"\n"}}
          {{"\t "}}<input type="text" class="form-control" id="recipient-{{$k}}" name="{{$k}}" value = "New" > {{"\n"}}
		{{- end -}}
   {{ end }}

          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Выйти</button>
        <button type="button" class="btn btn-primary" id="SaveInsertModal" >Записать</button>
      </div>
    </div>
  </div>
</div>

<!-- EDIT RECORDS -->
<div class="modal fade" id="editModal" tabindex="-1" role="dialog" aria-labelledby="editModalLabel" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="editModalLabel">{{$.TitleFormEdit}}Редактирование данных</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form id="ajaxEditForm" method="POST" action='medit' >
          <div class="form-group">
   {{ range $_, $v := $.Tabheader }}
		{{- range $k, $v := $v -}}
		{{"\t "}}<label for="edit-{{$k}}" class="col-form-label">{{$v}}:</label>{{"\n"}}
        {{"\t "}}<input type="text" class="form-control" id="edit-{{$k}}" name="{{$k}}" disabled>{{"\n"}}
		{{- end -}}
   {{ end }}
         </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Выйти</button>
        <button type="button" class="btn btn-primary" id="SaveEditModal" >Записать</button>
      </div>
<!-- ДЛЯ ОТЛАДКИ просмотра ЗАПРОСОВ
  <div class="container">
    <div class="well col-xs-12">
      <div class="control-label col-xs-12">
        <label>Data sent:</label>
      </div>
      <div class="col-xs-12">
        <textarea class="form-control" readonly id="dataSent">here: data sent...</textarea>
        <br>
      </div>
      <div class="control-label col-xs-12">
        <label>Result:</label>
      </div>
      <div class="col-xs-12">
        <textarea class="form-control" readonly id="results">Waiting to send request</textarea>
      </div>
    </div>
  </div>
-->
    </div>
  </div>
</div>

</body>

<script>
$(document).ready(function() {

  var table = $('#myTable').DataTable({
        "language": {
            "url": "static/DataTables/Russian.json"
        },
          "lengthMenu":[[16,24,50, -1], [16,24, 50, "All"]],
    dom: 'B<"clear">lfrtip',
   keys: {
        columns: ':not(:last-child)'
    },
    buttons: true,

    buttons: [
        {
                extend: 'csv',
                text: 'Отчет CSV',
                exportOptions: {
                        columns: {{ $.Columns }}
                },
                customize: function (csv) {
                        var split_csv = csv.split("\n");
                        split_csv[0] = '{{$.Hcsv}}';

                        $.each(split_csv.slice(1), function (index, csv_row) {
                                var csv_cell_array = csv_row.split('","');
                                csv_cell_array[0] = csv_cell_array[0].replace(/"/g, '');
                                csv_cell_array[3] = csv_cell_array[3].replace(/"/g, '');
                                csv_cell_array_quotes = '"' + csv_cell_array.join('";"') + '"';
                                split_csv[index + 1] = csv_cell_array_quotes;
                        });
                        csv = split_csv.join("\n");
                        return csv;
                }
        },

        {
                extend: 'excel',
                messageTop: 'Мониторинг версий ПО МАС ОПС УФПС Липецкой обл.',
                filename: 'file_excel',
                text: 'Отчет EXCEL',
                exportOptions: {
                        columns: {{ $.Columns }}
                }
        },
        {
                extend: 'pdf',
                filename: 'file_name',
                text: 'Отчет PDF',
                messageTop: 'Мониторинг версий ПО МАС ОПС УФПС Липецкой обл.',
                exportOptions: {
                        columns: {{ $.Columns}}
                }
        }

]
  }); 

<!-- Открыть модальную форму добавления записи -->
$('#insertModal').on('show.bs.modal', function (event) {
});

<!--Передать на добавление в БД новой записи -->
$("#SaveInsertModal").click(function(event) {
    event.preventDefault();
    var form = $('#ajaxCreateForm');
    var method = form.attr('method');
    var url = form.attr('action'); <!-- mcreate -->
    var formdata = form.serialize();
    console.log(formdata);
<!--    ajaxCallRequest(method, url, formdata); -->
    if (method =='POST') {
    $.post('mcreate',formdata, processData);
        function processData(data){
                if (data == 'pass') { 
                 console.log('Запись успешно добавлена '+formdata);
                 document.location.href = '/'                                       
                } else {                        
                 alert('Ошибка записи !!! '+ data);                                        
                }                                                
        }  
     } 
  });

<!-- Выбрана ссылка "Редактировать" на записи -->
$('.edit_modal').click(function(event) {
 	event.preventDefault();
 	var $editRow =null;
<!--    получим значения из таблици -->
 	$editRow = $(event.target ).closest( "tr" );
 	$id = $editRow.data('tr-id');
    {{ range $_, $z := $.Tabtd}}
	{{- range $k, $v := $z -}}
          {{"\t"}}{{$v}} = $editRow.children('td.{{$v}}').text().trim();{{"\n"}}
	{{- end -}}
    {{ end }}

    console.log('Edit id='+$id);
    var $editForm = $('#ajaxEditForm');
        $editForm.find('#edit-id').val($id);
    {{ range $_, $z := $.Tabtd }}
	{{- range $k, $v := $z -}}
           {{"\t"}} $editForm.find('#edit-{{$v}}').val({{$v}});{{"\n"}}
	{{- end -}}
    {{ end }}
}); <!--.edit_modal -->


<!--Передать данные серверу для записи в БД -->
$("#SaveEditModal").click(function(event) {
    event.preventDefault();
    var form = $('#ajaxEditForm');
    var method = form.attr('method');
    var url = form.attr('action'); <!-- medit -->
<!-- var formdata = $(form).serialize(); -->
<!-- включаем в список на сериализацию в том числе и поля с атрибутом disabled -->
    var formdata = form.serializeIncludeDisabled();
    console.log(formdata);
<!--   Для ОТЛАДКИ и ПРОСМОТРА ЗАПРОСОВ -->
<!-- ajaxCallRequest(method, url, formdata);
    if (method =='POST') {
    $.post('medit',formdata, processData);
        function processData(data){
                if (data == 'pass') {                            
                 console.log('Запись успешно обновлена '+formdata);
                 document.location.href = '/'                                       
                } else {                        
                 alert('Ошибка записи на сервер ! '+ data);
                }                                                
        }  
     } 
  });

<!--Функция сериализации полей формы в том числе и с атрибутом disabled -->
$.fn.serializeIncludeDisabled = function () {
    let disabled = this.find(":input:disabled").removeAttr("disabled");
    let serialized = this.serialize();
    disabled.attr("disabled", "disabled");
    return serialized;
};

<!--Универсальная Функция faking ajax requests -->
function ajaxCallRequest(f_method, f_url, f_data) {
    $("#dataSent").val(unescape(f_data));
    var f_contentType = 'application/x-www-form-urlencoded; charset=UTF-8';
    $.ajax({
      url: f_url,
      type: f_method,
      contentType: f_contentType,
      dataType: 'json',
      data: f_data,
      success: function(data) {
        var jsonResult = JSON.stringify(data);
        $("#results").val(unescape(jsonResult));
      }
    });
  }

}); 
</script>
</html>
`
type kv struct {
     k,v string
}

	m := make(map[string]interface{})
	if err = json.Unmarshal([]byte(jsonStr), &m); err != nil {
		fmt.Println("Error Unmarshal")
		panic(err)
	}

	tf := template.FuncMap{
		"isInt": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return true
			default:
				return false
			}
		},
		"isString": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.String:
				return true
			default:
				return false
			}
		},
		"isSlice": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Slice:
				return true
			default:
				return false
			}
		},
		"isArray": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Array:
				return true
			default:
				return false
			}
		},
		"isMap": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			switch v.Kind() {
			case reflect.Map:
				return true
			default:
				return false
			}

		},

		"first": func (n int, sk []string, sv []string) []string { 
			return sv[0:n]
                },
}
	

	t := template.New("hello").Funcs(tf)
	tt, err := t.Parse(text1)
	if err != nil {
		fmt.Println("Error Parse text1")
		panic(err)
	}
	if err = tt.Execute(os.Stdout, &m); err != nil {
		fmt.Println("Error Execute text1")
		panic(err)
	}

	fmt.Fprint(os.Stdout, text2)

	tt3, err := t.Parse(text3)
	if err != nil {
		fmt.Println("Error Parse text3")
		panic(err)
	}
	if err = tt3.Execute(os.Stdout, &m); err != nil {
		fmt.Println("Error Execute text3")
		panic(err)
	}
	fmt.Fprint(os.Stdout, text4)

	tt5, err := t.Parse(text5)
	if err != nil {
		fmt.Println("Error Parse text5")
		panic(err)
	}
	if err = tt5.Execute(os.Stdout, &m); err != nil {
		fmt.Println("Error Execute text5")
		panic(err)
	}

}