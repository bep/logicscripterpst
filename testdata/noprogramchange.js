// A script that drops all program change events.
function HandleMIDI(event) {
	if (!(event instanceof ProgramChange)) {
		event.send();
	}
}
