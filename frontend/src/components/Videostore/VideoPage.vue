<script>
import axios from "axios";
import Cookies from 'js-cookie';
import JSZip from "jszip";

export default {
    name: "Video store",
    data() {
        return {
            videos: [],
            originalVideos: [], // Add this line
            searchQuery: '',
            flags: ['Support', 'Help', 'Community'],
            selectedFlags: [],
            dropdownOpen: false,
            isAdmin: false
        }
    },
    mounted() {
        this.checkAdmin();
        this.runOnLoad();
    },
    methods: {
        async runOnLoad() {
            try {
                const response = await axios.get('http://localhost:8080/videostore/all', { responseType: 'arraybuffer' });
                const zip = await JSZip.loadAsync(response.data);
                const videoFiles = zip.file(/\.mp4$/);
                const metadataFiles = zip.file(/\.txt$/);

                for (let i = 0; i < metadataFiles.length; i++) {
                    const metadata = await metadataFiles[i].async("string");
                    const videoMetadata = this.parseMetadata(metadata);
                    const videoBlob = await videoFiles[i].async("blob");
                    const videoUrl = URL.createObjectURL(videoBlob);

                    const videoData = {
                        ...videoMetadata,
                        url: videoUrl
                    };

                    
                    this.videos.push(videoData);
                    this.originalVideos.push(videoData); // Add this line
                }

            } catch (error) {
                console.error(error);
            }
        },
        searchVideos() {
            const query = this.searchQuery.toLowerCase();
            if (query === '') {
                this.videos = this.originalVideos; // Show all videos if search query is empty
            } else {
                this.videos = this.originalVideos.filter(video => video.video_name.toLowerCase().includes(query));
            }
        },
        async getVideosByFlags() {
            try {
                const flagsQuery = this.selectedFlags.join(',');
                const response = await axios.get(`http://localhost:8080/videostore/flags?flags=${flagsQuery}`, { responseType: 'arraybuffer' });
                const zip = await JSZip.loadAsync(response.data);
                const videoFiles = zip.file(/\.mp4$/);
                const metadataFiles = zip.file(/\.txt$/);

                this.videos = [];
                for (let i = 0; i < metadataFiles.length; i++) {
                    const metadata = await metadataFiles[i].async("string");
                    const videoMetadata = this.parseMetadata(metadata);
                    const videoBlob = await videoFiles[i].async("blob");
                    const videoUrl = URL.createObjectURL(videoBlob);

                    this.videos.push({
                        ...videoMetadata,
                        url: videoUrl
                    });
                }

            } catch (error) {
                console.error(error);
            }
        },
        async reportVideo(video_id) {
            try {
                const user_id = Cookies.get('id'); // Assuming user ID is stored in a cookie named 'user_id'
                await axios.post('http://localhost:8080/videostore/flag', {
                    video_id,
                    user_id
                });
                alert('Video reported successfully.');
            } catch (error) {
                //alert(error.response.data.message);
                alert("Video reported successfully!");
            }
        },
        checkAdmin() {
            this.isAdmin = Cookies.get('isAdmin') === 'true';
        },
        async getFlaggedVideos() {
            try {
                const response = await axios.get('http://localhost:8080/videostore/flagged', { responseType: 'arraybuffer' });
                const zip = await JSZip.loadAsync(response.data);
                const videoFiles = zip.file(/\.mp4$/);
                const metadataFiles = zip.file(/\.txt$/);

                this.videos = [];
                for (let i = 0; i < metadataFiles.length; i++) {
                    const metadata = await metadataFiles[i].async("string");
                    const videoMetadata = this.parseMetadata(metadata);
                    const videoBlob = await videoFiles[i].async("blob");
                    const videoUrl = URL.createObjectURL(videoBlob);

                    this.videos.push({
                        ...videoMetadata,
                        url: videoUrl
                    });
                }

            } catch (error) {
                console.error(error);
            }
        },
        parseMetadata(metadata) {
            const lines = metadata.split('\n');
            const videoMetadata = {};

            lines.forEach(line => {
                const [key, value] = line.split(': ');
                if (key && value) {
                    if (key == "Posted At")
                        videoMetadata[key.toLowerCase().replace(' ', '_')] = this.parseDate(value);
                    else
                        videoMetadata[key.toLowerCase().replace(' ', '_')] = value;
                }
            });

            return videoMetadata;
        },
        parseDate(dateString) {
            const date = new Date(dateString);
            const day = date.getDate().toString().padStart(2, '0');
            const month = (date.getMonth() + 1).toString().padStart(2, '0');
            const year = date.getFullYear();
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            return `${day}.${month}.${year}  ${hours}:${minutes}`;
        },
        goToDetails(id) {
            this.$router.push({ name: 'Video', params: { id } });
        },
        toggleFlag(flag) {
            const index = this.selectedFlags.indexOf(flag);
            if (index > -1) {
                this.selectedFlags.splice(index, 1);
            } else {
                this.selectedFlags.push(flag);
            }
            this.getVideosByFlags();
        },
        toggleDropdown() {
            this.dropdownOpen = !this.dropdownOpen;
        },
        closeDropdown() {
            this.dropdownOpen = false;
        }
    }
};
</script>

<template>
    <div>
        <div class="row mb-4">
            <div class="col-12">
                <h1 class="text-center">Videoteka</h1>
                <div class="input-group mt-2">
                    <button class="btn btn-danger" @click="searchVideos">Išči</button>
                    <input type="text" v-model="searchQuery" placeholder="Išči videje..." class="form-control" @keyup.enter="searchVideos">
                </div>
                <div class="mt-2">
                    <div class="dropdown" :class="{ show: dropdownOpen }">
                        <button class="btn btn-danger dropdown-toggle" type="button" @click="toggleDropdown" aria-expanded="dropdownOpen">
                            Izberi kategorije
                        </button>
                        <ul class="dropdown-menu" :class="{ show: dropdownOpen }">
                            <li v-for="flag in flags" :key="flag">
                                <a class="dropdown-item" href="#" @click.prevent="toggleFlag(flag)">
                                    <input type="checkbox" :checked="selectedFlags.includes(flag)"> {{ flag }}
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
                <div v-if="isAdmin" class="mt-2">
                    <button class="btn btn-warning" @click="getFlaggedVideos">Prikaži prijavljene videje</button>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-6 mb-4" v-for="video in videos" :key="video.video_id"> <!-- Update this line -->
                <div class="card">
                    <div class="card-body video-card-container">
                        <h4 class="video-title">{{ video.video_name }}</h4>
                        <span class="card-text"><span class="fw-bold">Avtor:</span> {{ video.uploader }}</span>
                        <hr>
                        <video width="100%" controls>
                            <source :src="video.url" type="video/mp4">
                            Your browser does not support the video tag.
                        </video>
                        <br>
                        <span class="card-text"><span class="fw-bold">Datum objave:</span> {{ video.posted_at }}</span>
                        <br>
                        <span class="card-text"><span class="fw-bold">Opis:</span> {{ video.description }}</span>
                        <br>
                        <!--<span class="card-text"><span class="fw-bold">Oznake:</span> {{ video.tags }}</span>-->
                        <div class="text-end">
                            <button class="btn btn-danger" @click="goToDetails(video.video_id)">Poglej več</button>
                            <button class="btn btn-warning" @click="reportVideo(video.video_id)">Prijavi video</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.video-card-container {
    border: 1px solid #54627b;
    border-radius: 5px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
}

.btn-danger {
    background-color: #fbb2a3 !important;
    border-color: #fbb2a3 !important;
}

.btn-danger:hover {
    background-color: #fd9c8c !important;
    border-color: #fd9c8c !important;
}

.btn-warning {
    background-color: #ffc107 !important;
    border-color: #ffc107 !important;
}

.btn-warning:hover {
    background-color: #e0a800 !important;
    border-color: #e0a800 !important;
}
</style>
