document.getElementById('init-form').addEventListener('submit', function(e) {
    e.preventDefault();

    const num_frames = document.getElementById('num_frames').value;
    const algorithm = document.getElementById('algorithm').value;
    const page_sequence = document.getElementById('page_sequence').value.split(',').map(Number);

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
        } else {
            alert('Error initializing simulation.');
        }
    });
});

document.getElementById('next-step').addEventListener('click', function() {
    const sim_id = document.getElementById('simulation_id').innerText;
    fetch(`/api/simulation/${sim_id}/next`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        updateState(sim_id);
        alert(data.action);
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
    });
}
