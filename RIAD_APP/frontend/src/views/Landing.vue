<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const rooms = [
  { name: 'Chambre Supérieure', price: '120', image: '/images/image2.jpg', desc: 'Chambre élégante avec vue sur le patio intérieur' },
  { name: 'Suite Al Sagia', price: '180', image: '/images/image3.jpg', desc: 'Suite spacieuse avec salon privé et terrasse' },
  { name: 'Suite Familiale', price: '250', image: '/images/image4.jpg', desc: 'Deux chambres communicantes pour famille' },
]

const gallery = [
  '/images/images.jpg', '/images/image5.jpg',
  '/images/image2.jpg', '/images/image4.jpg',
  '/images/image3.jpg', '/images/image5.jpg',
]
</script>

<template>
  <div class="min-h-screen bg-[#F4F1EA]">
    <!-- Nav -->
    <nav class="bg-[#2D4635] text-white py-3 sm:py-4 px-4 sm:px-8 flex items-center justify-between sticky top-0 z-50">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-white/10 flex items-center justify-center">
          <span class="text-sm">🌙</span>
        </div>
        <span class="font-serif text-base sm:text-xl tracking-widest">RIAD AL SAGIA</span>
      </div>
      <div class="flex items-center gap-3 sm:gap-4">
        <button @click="router.push('/login')"
          class="text-[10px] sm:text-xs uppercase tracking-widest hover:opacity-70 transition hidden sm:block">
          Connexion
        </button>
        <button @click="router.push('/register')"
          class="bg-white text-[#2D4635] px-4 sm:px-6 py-2 sm:py-2.5 text-[10px] sm:text-xs font-bold uppercase tracking-tighter hover:bg-gray-100 transition rounded-sm">
          S'inscrire
        </button>
      </div>
    </nav>

    <!-- Hero -->
    <header class="relative h-[65vh] sm:h-[80vh] flex items-center justify-center text-center text-white overflow-hidden">
      <div class="absolute inset-0">
        <div class="absolute inset-0 bg-gradient-to-b from-black/40 via-black/30 to-black/60 z-10"></div>
        <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1549294413-26f195af03c1?auto=format&fit=crop&q=80')] bg-cover bg-center scale-105"></div>
      </div>
      <div class="relative z-20 px-4 max-w-3xl animate-in">
        <span class="uppercase tracking-[0.3em] sm:tracking-[0.4em] text-[10px] sm:text-sm mb-2 sm:mb-3 block opacity-80 font-light">
          Riad Al Sagia — Marrakech
        </span>
        <h1 class="text-4xl sm:text-6xl md:text-7xl font-serif mb-3 sm:mb-4 leading-tight">L'Art de l'Accueil</h1>
        <p class="text-sm sm:text-lg md:text-xl font-light italic mb-6 sm:mb-8 opacity-90 px-4">
          Une expérience marocaine authentique au cœur de la Médina
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4">
          <button @click="router.push('/register')"
            class="w-full sm:w-auto bg-white text-[#2D4635] px-8 sm:px-10 py-3 uppercase text-xs tracking-widest font-bold hover:bg-gray-100 transition shadow-lg">
            Réserver maintenant
          </button>
          <button @click="router.push('/login')"
            class="w-full sm:w-auto border border-white/60 text-white/90 px-8 sm:px-10 py-3 uppercase text-xs tracking-widest hover:bg-white/10 transition">
            Se connecter
          </button>
        </div>
      </div>
    </header>

    <!-- Stats bar -->
    <section class="bg-white border-b border-gray-100">
      <div class="max-w-5xl mx-auto grid grid-cols-2 sm:grid-cols-4 divide-x divide-gray-100">
        <div class="py-6 sm:py-8 text-center">
          <p class="font-serif text-2xl sm:text-3xl text-[#2D4635]">8</p>
          <p class="text-[10px] sm:text-xs text-gray-400 uppercase tracking-widest mt-1">Chambres</p>
        </div>
        <div class="py-6 sm:py-8 text-center">
          <p class="font-serif text-2xl sm:text-3xl text-[#2D4635]">15</p>
          <p class="text-[10px] sm:text-xs text-gray-400 uppercase tracking-widest mt-1">Années d'histoire</p>
        </div>
        <div class="py-6 sm:py-8 text-center">
          <p class="font-serif text-2xl sm:text-3xl text-[#2D4635]">4.9</p>
          <p class="text-[10px] sm:text-xs text-gray-400 uppercase tracking-widest mt-1">Note moyenne</p>
        </div>
        <div class="py-6 sm:py-8 text-center">
          <p class="font-serif text-2xl sm:text-3xl text-[#2D4635]">98%</p>
          <p class="text-[10px] sm:text-xs text-gray-400 uppercase tracking-widest mt-1">Satisfaction</p>
        </div>
      </div>
    </section>

    <!-- About -->
    <section class="py-16 sm:py-20 max-w-5xl mx-auto px-6 grid grid-cols-1 md:grid-cols-2 gap-8 sm:gap-12 items-center">
      <div class="aspect-[4/5] overflow-hidden rounded-sm">
        <div class="w-full h-full bg-[url('/images/images.jpg')] bg-cover bg-center hover:scale-105 transition duration-700"></div>
      </div>
      <div>
        <span class="text-[10px] sm:text-xs uppercase tracking-widest text-[#2D4635] font-bold">Notre Histoire</span>
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-serif text-[#2D4635] mt-3 mb-6 leading-tight">
          Un Havre de Paix<br />Historique
        </h2>
        <div class="w-12 h-[1px] bg-[#2D4635] mb-6"></div>
        <p class="text-gray-600 leading-relaxed text-sm sm:text-base font-light">
          Niché dans un quartier historique de Marrakech, le Riad Al Sagia vous accueille dans un
          cadre alliant tradition marocaine et confort contemporain. Restauré avec soin, chaque
          recoin raconte une histoire.
        </p>
        <p class="text-gray-500 leading-relaxed text-sm sm:text-base font-light mt-4">
          Un service personnalisé vous attend pour un séjour inoubliable au cœur de la Médina.
        </p>
        <button @click="router.push('/register')"
          class="mt-6 bg-[#2D4635] text-white px-6 py-2.5 text-xs uppercase tracking-widest hover:bg-[#3a5a45] transition rounded-sm">
          En savoir plus
        </button>
      </div>
    </section>

    <!-- Rooms -->
    <section class="py-16 sm:py-20 bg-white">
      <div class="text-center mb-10 sm:mb-12 px-4">
        <span class="text-[10px] sm:text-xs uppercase tracking-widest text-gray-400 font-bold">Séjournez</span>
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-serif text-[#2D4635] mt-2 mb-3">Nos Chambres & Suites</h2>
        <p class="text-gray-400 text-sm italic max-w-md mx-auto">Havre de paix au cœur de la Médina</p>
      </div>
      <div class="max-w-6xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8 px-4 sm:px-8">
        <div v-for="room in rooms" :key="room.name"
          class="group cursor-pointer" @click="router.push('/register')">
          <div class="overflow-hidden mb-4 aspect-[4/5] bg-gray-100">
            <div class="w-full h-full bg-cover bg-center group-hover:scale-105 transition duration-700"
              :style="{ backgroundImage: `url(${room.image})` }"></div>
          </div>
          <div class="px-1">
            <div class="flex items-center justify-between mb-1">
              <h3 class="font-serif text-[#2D4635] text-lg sm:text-xl">{{ room.name }}</h3>
              <span class="text-xs text-gray-400 uppercase tracking-tighter">{{ room.price }} €</span>
            </div>
            <p class="text-xs sm:text-sm text-gray-500 font-light">{{ room.desc }}</p>
          </div>
        </div>
      </div>
      <div class="text-center mt-10 sm:mt-12">
        <button @click="router.push('/register')"
          class="bg-[#2D4635] text-white px-8 sm:px-10 py-3 text-xs uppercase tracking-widest hover:bg-[#3a5a45] transition rounded-sm">
          Voir toutes les chambres
        </button>
      </div>
    </section>

    <!-- Gallery -->
    <section class="py-16 sm:py-20 px-4 sm:px-8 max-w-7xl mx-auto">
      <div class="text-center mb-10 sm:mb-12">
        <span class="text-[10px] sm:text-xs uppercase tracking-widest text-gray-400 font-bold">Galerie</span>
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-serif text-[#2D4635] mt-2 mb-3">Notre Riad en Images</h2>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-3 sm:gap-4">
        <div v-for="(img, i) in gallery" :key="i"
          class="overflow-hidden rounded-sm aspect-[4/3] bg-gray-100"
          :class="i === 0 ? 'md:col-span-2 md:row-span-2' : ''">
          <div class="w-full h-full bg-cover bg-center hover:scale-105 transition duration-700"
            :style="{ backgroundImage: `url(${img})` }"></div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="py-16 sm:py-20 bg-[#2D4635] text-white">
      <div class="max-w-6xl mx-auto px-6 text-center">
        <span class="text-[10px] sm:text-xs uppercase tracking-widest text-[#e2bc8a] font-bold">Services</span>
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-serif mt-2 mb-12">Une Expérience Complète</h2>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-8 sm:gap-12">
          <div class="group">
            <div class="w-14 h-14 mx-auto mb-5 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7 text-[#e2bc8a]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z"/></svg>
            </div>
            <h3 class="font-serif text-lg sm:text-xl mb-2">Spa & Hammam</h3>
            <p class="text-white/60 text-sm font-light">Détendez-vous dans notre hammam traditionnel après une journée dans la Médina</p>
          </div>
          <div class="group">
            <div class="w-14 h-14 mx-auto mb-5 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7 text-[#e2bc8a]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42"/></svg>
            </div>
            <h3 class="font-serif text-lg sm:text-xl mb-2">Restaurant</h3>
            <p class="text-white/60 text-sm font-light">Cuisine marocaine authentique préparée avec des produits du terroir</p>
          </div>
          <div class="group">
            <div class="w-14 h-14 mx-auto mb-5 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7 text-[#e2bc8a]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 9v.906a2.25 2.25 0 01-1.183 1.981l-6.478 3.488M2.25 9v.906a2.25 2.25 0 001.183 1.981l6.478 3.488m8.839 2.51l-4.66-2.51m0 0l-1.023-.55a2.25 2.25 0 00-2.134 0l-1.022.55m0 0l-4.661 2.51m16.5 1.615a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V8.844a2.25 2.25 0 011.183-1.981l7.5-4.039a2.25 2.25 0 012.134 0l7.5 4.039a2.25 2.25 0 011.183 1.98V19.5z"/></svg>
            </div>
            <h3 class="font-serif text-lg sm:text-xl mb-2">Conciergerie</h3>
            <p class="text-white/60 text-sm font-light">Service personnalisé pour organiser vos excursions et transferts</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Testimonial -->
    <section class="py-16 sm:py-20 bg-white text-center px-6">
      <div class="max-w-2xl mx-auto">
        <span class="text-4xl sm:text-5xl text-[#2D4635]/20 font-serif leading-none">"</span>
        <p class="text-lg sm:text-xl text-gray-600 font-light italic mt-2 leading-relaxed">
          Un véritable joyau caché. L'accueil, le cadre, la cuisine... tout était parfait.
          Une expérience marocaine authentique que je recommande vivement.
        </p>
        <div class="mt-6">
          <p class="font-serif text-[#2D4635] text-sm">— Sarah & Thomas</p>
          <p class="text-[10px] sm:text-xs text-gray-400 uppercase tracking-widest mt-1">Voyageurs, France</p>
        </div>
        <div class="flex items-center justify-center gap-1 mt-4">
          <span v-for="i in 5" :key="i" class="text-[#e2bc8a] text-sm">★</span>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="py-16 sm:py-20 bg-gradient-to-br from-[#2D4635] to-[#1a2f20] text-center text-white px-6">
      <h2 class="text-2xl sm:text-3xl md:text-4xl font-serif mb-3">Prêt à Vivre l'Expérience ?</h2>
      <p class="text-sm sm:text-base text-white/70 mb-8 font-light max-w-md mx-auto">
        Réservez votre séjour et laissez-vous envoûter par la magie de Marrakech
      </p>
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4">
        <button @click="router.push('/register')"
          class="w-full sm:w-auto bg-white text-[#2D4635] px-8 sm:px-10 py-3 text-xs uppercase tracking-widest font-bold hover:bg-gray-100 transition rounded-sm shadow-lg">
          Créer un compte
        </button>
        <button @click="router.push('/login')"
          class="w-full sm:w-auto border border-white/40 text-white/80 px-8 sm:px-10 py-3 text-xs uppercase tracking-widest hover:bg-white/10 transition rounded-sm">
          Déjà membre
        </button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-[#1a2f20] text-[#e2bc8a] py-12 sm:py-16 px-6 sm:px-10">
      <div class="max-w-7xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 sm:gap-12">
        <div>
          <div class="flex items-center gap-2 mb-4">
            <span class="text-lg">🌙</span>
            <h4 class="font-serif text-lg">Riad Al Sagia</h4>
          </div>
          <p class="text-xs leading-relaxed opacity-60">
            Situé à quelques minutes de la place Jemaa el-Fna, dans un quartier sécurisé de la Médina.
          </p>
        </div>
        <div>
          <h5 class="text-xs uppercase tracking-widest mb-4 font-bold text-white">Liens</h5>
          <ul class="text-xs space-y-2 opacity-60">
            <li><a href="#" class="hover:opacity-100 transition">Accueil</a></li>
            <li><a href="#" class="hover:opacity-100 transition">Chambres</a></li>
            <li><a href="#" class="hover:opacity-100 transition">Services</a></li>
            <li><a href="#" class="hover:opacity-100 transition">Contact</a></li>
          </ul>
        </div>
        <div>
          <h5 class="text-xs uppercase tracking-widest mb-4 font-bold text-white">Contact</h5>
          <p class="text-xs opacity-60 leading-relaxed">
            Place de la Bahia, Médina<br />
            Marrakech 40000, Maroc<br />
            +212 524 XX XX XX<br />
            info@riadalsagia.ma
          </p>
        </div>
        <div>
          <h5 class="text-xs uppercase tracking-widest mb-4 font-bold text-white">Application</h5>
          <p class="text-xs opacity-60 mb-3">Gérez votre séjour depuis l'application Riad Manager</p>
          <button @click="router.push('/login')"
            class="text-xs uppercase tracking-widest border border-[#e2bc8a]/40 px-4 py-2 hover:bg-white/5 transition rounded-sm">
            Se connecter
          </button>
        </div>
      </div>
      <div class="max-w-7xl mx-auto mt-10 pt-6 border-t border-white/5 text-center text-xs opacity-40">
        &copy; 2026 Riad Al Sagia. Tous droits réservés.
      </div>
    </footer>
  </div>
</template>

<style scoped>
.font-serif {
  font-family: 'Playfair Display', Georgia, serif;
}
</style>
