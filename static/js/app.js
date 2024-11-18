document.addEventListener("DOMContentLoaded", function() {
    // Fetch available algorithms
    fetch('/api/algorithms')
        .then(response => response.json())
        .then(data => {
            const algorithmSelect = document.getElementById('algorithm');
            data.algorithms.forEach(alg => {
                const option = document.createElement('option');
                option.value = alg;
                option.text = alg;
                algorithmSelect.add(option);
            });
        });

    // Handle form submission
    document.getElementById('simulationForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const numFrames = parseInt(document.getElementById('numFrames').value);
        const pages = document.getElementById('pages').value.split(',').map(p => parseInt(p.trim()));
        const vaInput = document.getElementById('virtualAddresses').value.split(',').map(va => va.trim());
        const virtualAddresses = vaInput.map(va => {
            const [page, offset] = va.split(':').map(v => parseInt(v.trim()));
            return { page, offset };
        });
        const algorithm = document.getElementById('algorithm').value;

        const payload = {
            numFrames,
            pages,
            virtualAddresses,
            algorithm
        };

        fetch('/api/simulate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })
        .then(response => response.json())
        .then(data => {
            const tbody = document.querySelector('#resultsTable tbody');
            tbody.innerHTML = '';
            data.simulationSteps.forEach(step => {
                const row = document.createElement('tr');

                const stepCell = document.createElement('td');
                stepCell.textContent = step.step;
                row.appendChild(stepCell);

                const pageCell = document.createElement('td');
                pageCell.textContent = step.page;
                row.appendChild(pageCell);

                const hitCell = document.createElement('td');
                hitCell.textContent = step.hit ? 'Hit' : 'Miss';
                row.appendChild(hitCell);

                const actionCell = document.createElement('td');
                actionCell.textContent = step.action;
                row.appendChild(actionCell);

                const paCell = document.createElement('td');
                paCell.textContent = step.physicalAddress;
                row.appendChild(paCell);

                tbody.appendChild(row);
            });

            document.getElementById('totalHits').textContent = data.totalHits;
            document.getElementById('totalMisses').textContent = data.totalMisses;
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred during simulation.');
        });
    });
});
