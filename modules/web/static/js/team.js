function create_team() {
    $.post('/team/create', {
        'name' : $("#name").val(),
        'resume' : $("#resume").val(),
        'users' : $("#users").val()
    }, function(json) {
        handle_json(json, function(){
            location.href = '/teams';
        });
    });
}

function edit_team(org_id, team_id) {
    $.post('/team/'+team_id+'/edit', {
        'resume' : $("#resume").val(),
        'users' : $("#users").val(),
    }, function(json) {
        handle_json(json, function(){
            location.href = '/teams';
        });
    });
}