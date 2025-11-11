import App from './App.svelte';
import { mount } from 'svelte';
import { addCollection } from '@iconify/svelte';
import circleFlagsData from '@iconify-json/circle-flags/icons.json';

addCollection(circleFlagsData);

const app = mount(App, {
  target: document.getElementById('app')!
});

export default app;
