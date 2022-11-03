<script>
	let matches = [];
	let match = [];
	let ct = 0;
	let t = 0;
	fetch(`/api/matches`).then((response) => {
    response.json().then((json) => {
      matches = json;
    });
  });


	function selectMatch(m) {

		fetch(`/api/matchjson/${m}`).then((response) => {
			 response.json().then((json) => {
       match = json;
    });
    });
		ct = 0;
		t = 0;
		match = m;
	}

function parseLine(line) {
		let roundnotice = "";
		if (!line) return "";
		console.log(line)
			if (line.type == "TeamScored") {
			console.log(line)
			if (line.side == "CT") {
				ct = line.score
			
			} else {
				t = line.score
			}
		}
		if (line.type == "TeamNotice") {
			return "--- Round over ---"
}
		if (line.type == "PlayerKill") {
			
			if (line.headshot) {
			return line.attacker.name + " killed " + line.victim.name + " with a flippin headshot with " + line.weapon

			} else {
				return line.attacker.name + " killed " + line.victim.name + " with " + line.weapon

			}
		} 
	
	return "";
	}
</script>

<main><div id="navigation">
	 {#each matches as match}
	 <div class="match" on:click={() => selectMatch(match)}>{match}</div>
	 {/each}
	 </div>
	 <div id="main">
	 	<pre>
	 		{#each match as line}
	 		<div>{parseLine(line)}</div>
	 		{/each}
	 		<div class="border">CT: {ct}, T: {t}</div>
	 </pre>
	 </div>
</main>

<style>
	#navigation {
	 width: 250px;
  position: absolute;
  top: 0;
  bottom: 0;
  background-color: green;
	}
	#navigation div {
		display: block;
	}
	#main {
		margin-left: 300px;
	}
	.border {
		border:  1px solid red;
		width:  100px;
	}
	.match {
		cursor: pointer;
	}
</style>