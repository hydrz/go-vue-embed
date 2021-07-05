<template>
	<h1>Hello App!</h1>
	<img src="@/assets/logo.png" alt="" />
	<p>
		<router-link to="/">Go to Home</router-link>
		|
		<router-link to="/about">Go to About</router-link>
	</p>
	<router-view></router-view>
</template>

<script setup>
// This starter template is using Vue 3 experimental <script setup> SFCs
// Check out https://github.com/vuejs/rfcs/blob/script-setup-2/active-rfcs/0000-script-setup.md
</script>

<style>
#app {
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	text-align: center;
	color: #2c3e50;
	margin-top: 60px;
}
</style>
