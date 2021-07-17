const searchForm = document.getElementById('search-form');
const searchField = document.getElementById('search-field');
const dialogBtn = document.getElementById('dialog-button');
const dialog = document.getElementById('dialog');
const cancelDialogBtn = document.getElementById('cancel-dialog');
const postForm = document.getElementById('post-form');

const toogleDialog = function toogleDialog() {
  dialog.classList.toggle('none');
};

const deleteBook = function deleteBook(id) {
  console.log(id);
};

const getDetail = function getDetail(id) {
  console.log(id);
};

const updateBook = function updateBook(id) {
  toogleDialog();
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

dialogBtn.addEventListener('click', () => {
  toogleDialog();
});

cancelDialogBtn.addEventListener('click', () => {
  toogleDialog();
});

searchForm.addEventListener('submit', (e) => {
  e.preventDefault();
  location.search = `name=${searchField.value}`;
});

postForm.addEventListener('submit', (e) => {
  e.preventDefault();
  postForm.reset();
});
