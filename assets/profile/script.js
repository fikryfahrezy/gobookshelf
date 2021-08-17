const form = document.getElementById('form');
const oldData = Object.fromEntries(new FormData(form).entries());

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const form = new FormData(e.target);
  const { email, name, address, password } = Object.fromEntries(form.entries());
  const data = {};

  if (email !== oldData.email) data.email = email;
  if (name !== oldData.name) data.name = name;
  if (address !== oldData.address) data.address = address;
  if (password !== oldData.password) data.password = password;

  fetch('/updateprofile', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'PATCH',
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (res.status >= 200 && res.status < 300) {
        window.location.reload();
      }
    })
    .catch((err) => {
      console.log(err);
    });
});
