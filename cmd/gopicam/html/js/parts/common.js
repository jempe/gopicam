const requestInit = {
	"cache": "no-store",
	"mode": "same-origin",
	"credentials": "same-origin",
	headers: {
		"Content-Type": "application/json; charset=utf-8"
	}
};

function handleResponse(response)
{
	if (document.querySelectorAll('button.disabled').length > 0)
	{
		document.querySelector("button.disabled").classList.remove("disabled");
	}

	if (response.status >= 200 && response.status < 300)
	{
		return Promise.resolve(response)
	} 
	else
	{
		return Promise.reject(new Error(response.statusText))
	}
}
function handleJson(response)
{
	return jsonResponse = response.json();
}

// display error messages
function log_error(message)
{
	console.log(message);
}

// get data of form
function get_form_data(form)
{
	let form_data = {};

	let fields = form.querySelectorAll("input, select");

	for(let field_index in fields)
	{
		let field = fields[field_index];

		if(field.name != "gorilla.csrf.Token")
		{
			form_data[field.name] = field.value;

		}
	}

	return form_data;
}

// get data of form
function populate_form_data(form, form_data, csrf_token)
{
	let fields = form.querySelectorAll("input, select");

	for(let field_index in fields)
	{
		let field = fields[field_index];

		if(field.name == "gorilla.csrf.Token")
		{
			field.value = csrf_token;

		}
		else if(form_data[field.name])
		{
			field.value = form_data[field.name];
		}
	}
}

