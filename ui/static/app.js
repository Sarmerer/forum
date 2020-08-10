$(document).ready(function() {
	$('#form').submit(function(e) {
		e.preventDefault();
		var form = $(this);
		var url = form.attr('action');
		$.ajax({
			type: 'POST',
			url: '/signin',
			data: form.serialize(), // serializes the form's elements.
			success: function(data) {
				alert(data); // show response from the php script.
			}
		});
	});
});
