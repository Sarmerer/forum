$(document).ready(function() {
	$('#form').submit(function(e) {
		e.preventDefault();
		var form = $(this);
		$.ajax({
			type: 'POST',
			url: '/signin',
			data: form.serialize(), // serializes the form's elements.
			success: function(data) {
				if (data.length > 0) {
				}
			}
		});
	});
	$('#signout').on('click', function() {
		$.ajax({
			type: 'POST',
			url: '/signout'
		});
	});
});
