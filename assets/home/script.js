const searchForm = document.getElementById('search-form');
const searchField = document.getElementById('search-field');
const dialogBtn = document.getElementById('dialog-button');
const dialog = document.getElementById('dialog');
const cancelDialogBtn = document.getElementById('cancel-dialog');
const submitBtn = document.getElementById('submit-btn');
const postForm = document.getElementById('post-form');

const smtBtnTxt = submitBtn.innerText;
let method = 'POST';
let bookId = '';

const toogleDialog = function toogleDialog() {
  dialog.classList.toggle('none');
};

const setFormField = function setFormField(id) {
  const name = document.getElementById(`name-${id}`).innerText;
  document.getElementById('form-name').value = name;

  const year = document.getElementById(`year-${id}`).innerText;
  document.getElementById('form-year').value = year;

  const author = document.getElementById(`author-${id}`).innerText;
  document.getElementById('form-author').value = author;

  const summary = document.getElementById(`summary-${id}`).innerText;
  document.getElementById('form-summary').value = summary;

  const publisher = document.getElementById(`publisher-${id}`).innerText;
  document.getElementById('form-publisher').value = publisher;

  const pageCount = document.getElementById(`pageCount-${id}`).innerText;
  document.getElementById('form-pageCount').value = pageCount;

  const readPage = document.getElementById(`readPage-${id}`).innerText;
  document.getElementById('form-readPage').value = readPage;

  const isReading = document.getElementById(`reading-${id}`).innerText;

  if (Boolean(Number(isReading))) {
    document.getElementById('form-true').checked = true;
  } else {
    document.getElementById('form-false').checked = true;
  }
};

const setFormButtonName = function setFormButtonName() {
  if (method.toUpperCase() === 'PUT') {
    submitBtn.innerText = 'Update';
  } else {
    submitBtn.innerText = smtBtnTxt;
  }
};

const deleteBook = function deleteBook(id) {
  fetch(`/books/${id}`, {
    method: 'DELETE',
  })
    .then((res) => res.json())
    .catch((err) => {
      console.log(err);
    })
    .finally(() => {
      location.reload();
    });
};

const updateBook = function updateBook(id) {
  toogleDialog();

  method = 'PUT';
  bookId = id;

  setFormField(id);
  setFormButtonName();
};

dialogBtn.addEventListener('click', () => {
  toogleDialog();

  method = 'POST';
  bookId = '';

  setFormButtonName();
});

cancelDialogBtn.addEventListener('click', () => {
  toogleDialog();

  method = 'POST';
  bookId = '';
});

searchForm.addEventListener('submit', (e) => {
  e.preventDefault();

  location.search = `name=${searchField.value}`;
});

postForm.addEventListener('submit', (e) => {
  e.preventDefault();

  const form = new FormData(e.target);
  const data = {
    name: form.get('name'),
    year: Number(form.get('year')),
    author: form.get('author'),
    summary: form.get('summary'),
    publisher: form.get('publisher'),
    pageCount: Number(form.get('pageCount')),
    readPage: Number(form.get('readPage')),
    reading: Boolean(Number(form.get('reading'))),
  };

  const url = bookId === '' ? '/books' : `/books/${bookId}`;
  const formMethod = method.toUpperCase();

  fetch(url, {
    method: formMethod,
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  })
    .then((res) => res.json())
    .catch((err) => {
      console.log(err);
    })
    .finally(() => {
      window.location.reload();
    });
});
