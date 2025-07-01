document.addEventListener('DOMContentLoaded', function () {
    const pdfOverlay = document.getElementById('pdfOverlay');
    const pdfViewer = document.getElementById('pdfViewer');
    const pdfClose = document.getElementById('pdfClose');
    const pdfLinks = document.querySelectorAll('.pdf-link');

    // Function to open the modal with the correct PDF
    const openModal = (pdfFile) => {
        if (pdfViewer && pdfOverlay) {
            const pdfPath = `/static/SeniorJury/${encodeURIComponent(pdfFile)}`;
            pdfViewer.setAttribute('src', pdfPath);
            pdfOverlay.classList.add('active');
            document.addEventListener('keydown', handleEscape);
        }
    };

    // Function to close the modal
    const closeModal = () => {
        if (pdfOverlay && pdfViewer) {
            pdfOverlay.classList.remove('active');
            pdfViewer.setAttribute('src', ''); // Stop PDF loading
            document.removeEventListener('keydown', handleEscape);
        }
    };

    // Add Escape key support for accessibility
    const handleEscape = (e) => {
        if (e.key === 'Escape' || e.key === 'Esc') {
            closeModal();
        }
    };

    // Attach click event to all "View Score" links
    pdfLinks.forEach(link => {
        link.addEventListener('click', function (e) {
            e.preventDefault();
            const filename = this.getAttribute('data-pdf');
            if (filename) {
                openModal(filename);
            }
        });
    });

    // Attach click event to the close button
    if (pdfClose) {
        pdfClose.addEventListener('click', closeModal);
    }

    // Attach click event to the overlay to close when clicking outside the PDF
    if (pdfOverlay) {
        pdfOverlay.addEventListener('click', function (e) {
            // Only close if the click is on the overlay itself, not the iframe
            if (e.target === pdfOverlay) {
                closeModal();
            }
        });
    }
});