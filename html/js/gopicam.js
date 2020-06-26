/*--include:parts/common.js:--*/
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
		if(response.status == 401)
		{
			show_login();
		}

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

		form_data[field.name] = field.value;
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

/*--includeend--*/

// show login form
function show_login()
{
	document.body.classList.add("show_login");
}

// hide login form
function hide_login()
{
	document.body.classList.remove("show_login");
}


//submit login form
function submit_login_form(form)
{
	document.querySelector(".login_form_container").classList.remove("error");

	let loginRequest = { 
		"cache": "no-store",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		"method" : "POST",
		"body" : "username=" + encodeURIComponent(document.getElementById("username").value) + "&password=" + encodeURIComponent(document.getElementById("password").value)
	}

	fetch("/api/login" , loginRequest).then(handleResponse).then(handleJson).then(function(data)
	{
		if(data.access != "granted")
		{
			//if can't login
			if(document.querySelectorAll("#login_form span.error").length == 0)
			{
				document.getElementById("login_form").insertAdjacentHTML('afterBegin', '<span class="error">Wrong username or password</span>');
			}

			document.querySelector(".login_form_container").classList.add("error");
		}
		else
		{
			hide_login();
		}

	}).catch(function(error)
	{
		log_error('Request failed' +  error);
	});

	return false;
}

function get_preview()
{
	let previewRequest = requestInit;
	previewRequest["method"] = "GET";

	fetch("/api/preview", previewRequest).then(handleResponse).then(handleJson).then(function(data)
	{

	}).catch(function(error)
	{
		log_error('Request failed' +  error);
	});
}

get_preview();
