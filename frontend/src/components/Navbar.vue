<template>
    <nav class="bg-gray-800 text-white">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <!-- Logo -->
                <div class="flex-shrink-0 flex items-center">
                    <FontAwesomeIcon :icon="faDice" class="mr-2" />
                    <router-link to="/" class="text-xl font-bold">ReRandom.neorg</router-link>
                </div>

                <!-- Desktop Menu -->
                <div class="hidden md:flex space-x-4 items-center">
                    <router-link
                        v-for="item in navItems"
                        :key="item.name"
                        :to="item.link"
                        class="hover:bg-gray-700 px-3 py-2 rounded-md text-sm font-medium"
                    >
                        {{ item.name }}
                    </router-link>
                </div>

                <!-- Mobile Menu Button -->
                <div class="md:hidden flex items-center">
                    <button
                        @click="toggleMenu"
                        class="bg-gray-800 inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                    >
                        <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path v-if="!isMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>
        </div>

        <!-- Mobile Menu -->
        <div v-if="isMenuOpen" class="md:hidden">
            <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-gray-800">
                <router-link
                    v-for="item in navItems"
                    :key="item.name"
                    :to="item.link"
                    class="hover:bg-gray-700 block px-3 py-2 rounded-md text-base font-medium"
                >
                    {{ item.name }}
                </router-link>
            </div>
        </div>
    </nav>
</template>

<script setup>
import { ref } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faDice } from '@fortawesome/free-solid-svg-icons';
  
    const navItems = ref([
    { name: 'List', link: { name: 'PredictsList' } },
    { name: 'Smart Generation', link: { name: 'SmartGeneration' } },
    { name: 'Login', link: { name: 'Login' } },
]);

const isMenuOpen = ref(false);

const toggleMenu = () => {
    isMenuOpen.value = !isMenuOpen.value;
};
</script>