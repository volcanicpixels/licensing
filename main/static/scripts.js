(function($){

  $.fn.serializeObject = function() {
    var o = {};
    var a = this.serializeArray();

    $.each(a, function() {
      if (o[this.name] !== undefined) {

        if (!o[this.name].push) {
          o[this.name] = [o[this.name]];
        }

        o[this.name].push(this.value || '');

      } else {
        o[this.name] = this.value || '';
      }
    });

    return o;
  };

  var serializeForm = function($section) {
    return $section.find('form').serializeObject();
  };

  var apiPOSTRequest = function(path, data) {
    return $.ajax({
      url: '/api' + path,
      type: 'POST',
      data: JSON.stringify(data),
      contentType: 'application/json; charset=utf-8',
      dataType: 'json',
    });
  };

  var setResult = function($section, data) {
    $section.find('.result').val(data.result);
  };

  var changeSection = function() {
    console.log(this);

    $('.section').removeClass('active');
    $('section').removeClass('active');

    var name = $(this).text().toLowerCase();

    $(this).addClass('active');
    $('.' + name).addClass('active');

  };

  $(document).ready(function(){
    var $create = $('.create');

    $('.section').click(changeSection);

    $('.create button').click(function(event){
      event.preventDefault();

      // disable button
      // start spinner

      apiPOSTRequest('/licenses', serializeForm($create))
        .done(function(data, textStatus){
          // show license
          console.log(data);
          setResult($create, data);
        })
        .fail(function(jqXHR){
          // Do something intelligent with the error
          alert(jqXHR.responseText);
          console.log(jqXHR.resultText);
        })
        .always(function(){
          // reset button and spinner
        });
    });

  });
}(jQuery));
