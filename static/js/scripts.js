document.getElementById('init-form').addEventListener('submit', function(e) {
    e.preventDefault();

    const num_frames = document.getElementById('num_frames').value;
    const algorithm = document.getElementById('algorithm').value;
    const page_sequence = document.getElementById('page_sequence').value.split(',').map(p => parseInt(p.trim())).filter(p => !isNaN(p));

    if (page_sequence.length === 0) {
        alert('Please enter a valid page sequence.');
        return;
    }

    fetch('/api/init', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            num_frames: parseInt(num_frames),
            algorithm: algorithm,
            page_sequence: page_sequence
        })
    })
    .then(response => response.json())
    .then(data => {
        if(data.simulation_id){
            document.getElementById('simulation_id').innerText = data.simulation_id;
            document.getElementById('init-section').style.display = 'none';
            document.getElementById('simulation-section').style.display = 'block';
            updateState(data.simulation_id);
            renderSteps(data.steps);
        } else {
            alert('Error initializing simulation: ' + (data.error || 'Unknown error.'));
        }
    })
    .catch(err => {
        console.error(err);
        alert('Error initializing simulation.');
    });
});

document.getElementById('next-step').addEventListener('click', function() {
    const sim_id = document.getElementById('simulation_id').innerText;
    fetch(`/api/simulation/${sim_id}/next`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        updateState(sim_id);
        appendStep(data.steps[data.steps.length - 1]);
        alert(data.action);
    })
    .catch(err => {
        console.error(err);
        alert('Error advancing simulation.');
    });
});

document.getElementById('reset-simulation').addEventListener('click', function() {
    const sim_id = document.getElementById('simulation_id').innerText;
    fetch(`/api/simulation/${sim_id}/reset`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        document.getElementById('current_step').innerText = data.current_step;
        document.getElementById('memory_frames').innerText = JSON.stringify(data.memory_frames);
        document.getElementById('page_faults').innerText = data.page_faults;
        document.getElementById('page_hits').innerText = data.page_hits;
        document.getElementById('status').innerText = data.status;
        clearSteps();
    })
    .catch(err => {
        console.error(err);
        alert('Error resetting simulation.');
    });
});

function updateState(sim_id){
    fetch(`/api/simulation/${sim_id}/state`)
    .then(response => response.json())
    .then(data => {
        document.getElementById('current_step').innerText = data.current_step;
        document.getElementById('memory_frames').innerText = JSON.stringify(data.memory_frames);
        document.getElementById('page_faults').innerText = data.page_faults;
        document.getElementById('page_hits').innerText = data.page_hits;
        document.getElementById('status').innerText = data.status;
    })
    .catch(err => {
        console.error(err);
        alert('Error fetching simulation state.');
    });
}

function renderSteps(steps){
    const tbody = document.querySelector('#steps-table tbody');
    tbody.innerHTML = ''; // Clear existing rows
    steps.forEach(step => {
        const row = document.createElement('tr');

        const stepNumberCell = document.createElement('td');
        stepNumberCell.innerText = step.step_number;
        row.appendChild(stepNumberCell);

        const virtualAddressCell = document.createElement('td');
        virtualAddressCell.innerText = step.virtual_address;
        row.appendChild(virtualAddressCell);

        const physicalFramesCell = document.createElement('td');
        physicalFramesCell.innerText = JSON.stringify(step.physical_frames);
        row.appendChild(physicalFramesCell);

        const actionCell = document.createElement('td');
        actionCell.innerText = step.action;
        row.appendChild(actionCell);

        tbody.appendChild(row);
    });
}

function appendStep(step){
    const tbody = document.querySelector('#steps-table tbody');
    const row = document.createElement('tr');

    const stepNumberCell = document.createElement('td');
    stepNumberCell.innerText = step.step_number;
    row.appendChild(stepNumberCell);

    const virtualAddressCell = document.createElement('td');
    virtualAddressCell.innerText = step.virtual_address;
    row.appendChild(virtualAddressCell);

    const physicalFramesCell = document.createElement('td');
    physicalFramesCell.innerText = JSON.stringify(step.physical_frames);
    row.appendChild(physicalFramesCell);

    const actionCell = document.createElement('td');
    actionCell.innerText = step.action;
    row.appendChild(actionCell);

    tbody.appendChild(row);
}

function clearSteps(){
    const tbody = document.querySelector('#steps-table tbody');
    tbody.innerHTML = '';
}
