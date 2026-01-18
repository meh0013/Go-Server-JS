submitFunc = async(e) => {
    const container = document.getElementById('container');
    if (!container.checkValidity()) {
      container.reportValidity();
      return;
    }
  
    let data = {
      // name: container.querySelector('input[name="name-input"]').value,
      name: document.getElementById("name-input").value,
      contact: document.getElementById("contact-input").value,
      website: document.getElementById("website-input").value,
      domain: document.getElementById("domain-input").value
  };
  
    try {
    let response = await fetch("http://localhost:8080/posts", {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data),
    });
      
    let text = await response.text();
      console.log(text);
      document.querySelector("#output").innerHTML = text;
      } 
      catch (err) {
      console.error(err);
      }
};