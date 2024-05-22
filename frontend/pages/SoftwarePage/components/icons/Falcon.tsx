import React from "react";

import type { SVGProps } from "react";

const Falcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    width="32"
    height="32"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <rect width="32" height="32" fill="#D31C12" />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M5.375 7.71218C5.375 7.77397 5.83297 8.70081 5.88749 8.74807C5.8984 8.75897 5.94565 8.83166 5.9929 8.91163C6.04015 8.99159 6.12011 9.097 6.16373 9.14788C6.21098 9.19877 6.24733 9.24965 6.24733 9.26056C6.24733 9.27146 6.28004 9.32235 6.32365 9.3696C6.33989 9.38878 6.3693 9.42295 6.40534 9.4648L6.40534 9.4648L6.40537 9.46484C6.45807 9.52605 6.52492 9.60368 6.58535 9.67491C7.05423 10.2201 7.73028 10.8526 8.44631 11.4123C8.46812 11.4268 8.52264 11.4704 8.56989 11.5104C8.61714 11.5504 8.66803 11.5831 8.67893 11.5831C8.68984 11.5831 8.74436 11.6195 8.79161 11.6631C8.92246 11.7794 9.05331 11.8703 9.74026 12.3282L10.0065 12.5049L10.0065 12.5049C10.3865 12.7572 10.4295 12.7858 10.6453 12.9134C10.7362 12.9679 10.8307 13.037 10.8561 13.0624C10.8816 13.0879 10.9215 13.1097 10.9397 13.1097C10.9652 13.1097 11.3323 13.3169 11.5176 13.4368C11.5467 13.4586 11.6121 13.495 11.663 13.5168C11.7939 13.5822 12.5208 14.022 12.7352 14.1674C13.3531 14.5817 14.1346 15.2905 14.3709 15.6431C14.469 15.7921 14.5344 15.8393 14.5344 15.7594C14.5344 15.7121 14.3854 15.4395 14.3272 15.3705C14.3018 15.345 14.28 15.305 14.28 15.2869C14.28 15.2469 14.0619 14.8907 13.9965 14.8216C13.971 14.7926 13.9202 14.7271 13.8874 14.6762C13.7893 14.5272 13.6694 14.3927 13.2986 14.0184C12.8516 13.5749 12.8007 13.535 12.1719 13.0661C11.8811 12.8516 11.6303 12.6626 11.6085 12.6445C11.5903 12.6299 11.5067 12.5754 11.4268 12.5281C11.3468 12.4809 11.2741 12.43 11.2632 12.4191C11.2523 12.4082 11.1796 12.3573 11.0996 12.3101C11.0197 12.2628 10.7907 12.1174 10.5908 11.9866C10.1873 11.7212 10.1183 11.6776 10.002 11.6158C9.95471 11.5904 9.91836 11.5577 9.91836 11.5395C9.91836 11.525 9.90019 11.5104 9.88202 11.5104C9.86021 11.5104 9.77298 11.4595 9.68938 11.3978C9.54763 11.296 9.11146 10.987 8.96244 10.8853C8.92246 10.8598 8.86067 10.8089 8.82069 10.7726C8.77707 10.7399 8.73345 10.7108 8.72255 10.7108C8.71165 10.7108 8.66076 10.6745 8.60987 10.6272C8.55899 10.5836 8.41724 10.4673 8.29366 10.3728C8.17008 10.2783 8.0465 10.1801 8.02105 10.1547C7.99561 10.1292 7.88294 10.0384 7.7739 9.95115C7.30502 9.58404 5.70576 8.0502 5.54583 7.81395C5.47677 7.71218 5.375 7.65039 5.375 7.71218ZM8.78436 8.37006C8.75528 8.39914 8.79526 8.52999 8.83161 8.52999C8.84978 8.52999 8.86432 8.54816 8.86432 8.5736C8.86432 8.62812 9.12602 9.15152 9.18417 9.21331C9.20598 9.23875 9.22779 9.2751 9.22779 9.29327C9.22779 9.34052 9.54037 9.79123 9.69303 9.96206C10.0565 10.3691 10.718 10.9834 10.9724 11.1506C11.0233 11.1833 11.1142 11.256 11.1723 11.3069C11.2341 11.3614 11.3468 11.4486 11.4268 11.5032C11.5067 11.5577 11.6049 11.6304 11.6412 11.6631C11.6812 11.6994 11.7793 11.7721 11.8593 11.8267C11.9429 11.8848 12.041 11.9575 12.0774 11.9902C12.1174 12.0266 12.2155 12.0993 12.2955 12.1538C12.3791 12.2083 12.4518 12.2628 12.4627 12.2737C12.4736 12.2846 12.5971 12.3828 12.7353 12.4918C12.877 12.6008 13.0079 12.7062 13.0297 12.7281C13.0515 12.7462 13.1969 12.8771 13.3532 13.0188C13.8838 13.495 14.2909 13.9384 14.5381 14.3091C14.778 14.6726 14.8725 14.8325 15.0869 15.2469C15.1996 15.4722 15.3668 15.7703 15.4504 15.912C15.5376 16.0538 15.6467 16.2355 15.6939 16.3191C15.7412 16.4027 15.8029 16.4936 15.832 16.519C15.8575 16.5445 15.8793 16.5845 15.8793 16.6026C15.8793 16.6208 15.8974 16.6572 15.9156 16.6826C15.9374 16.708 16.1083 16.9625 16.3009 17.2532C16.9152 18.1837 17.7584 18.9252 18.5871 19.256C18.9797 19.4123 19.8011 19.474 20.877 19.4268C21.6257 19.3977 22.4835 19.4195 22.6325 19.474C22.6798 19.4922 22.7852 19.5068 22.8688 19.5068C22.9524 19.5104 23.0869 19.5322 23.1668 19.5613C23.2468 19.5904 23.3449 19.6158 23.3849 19.6158C23.4249 19.6158 23.5049 19.6412 23.563 19.6703C23.6175 19.6994 23.6902 19.7248 23.7193 19.7248C23.7775 19.7248 23.9774 19.8121 24.439 20.0374C25.1041 20.3682 25.6421 20.7607 26.2091 21.3314C26.5871 21.713 26.6016 21.7239 26.6234 21.6294C26.6525 21.5167 26.6162 21.1315 26.5689 21.0442C26.5507 21.0079 26.4962 20.8952 26.449 20.7971C26.3581 20.5935 26.1873 20.35 25.9328 20.0556C25.8347 19.9356 25.7656 19.8193 25.7656 19.7685C25.7656 19.6194 25.6239 19.3868 25.4567 19.2632C25.2459 19.1069 24.5989 18.7798 24.5008 18.7798C24.4826 18.7798 24.4172 18.7544 24.3626 18.7253C24.3045 18.6962 24.2282 18.6708 24.1918 18.6708C24.1591 18.6708 24.1009 18.6562 24.0646 18.6381C24.0319 18.6199 23.8792 18.569 23.7302 18.5254C23.2323 18.3764 22.9124 18.2637 22.7997 18.1946C22.6798 18.1219 22.5671 17.9329 22.5671 17.8057C22.5671 17.7112 22.7052 17.2823 22.7488 17.2387C22.7706 17.2169 22.7852 17.1769 22.7852 17.1478C22.7852 17.1151 22.807 17.0461 22.8361 16.9915C22.8652 16.9407 22.9051 16.7807 22.9306 16.639C23.0505 15.9339 22.6689 15.4832 21.4404 14.8689C21.2768 14.7889 20.9787 14.6544 20.8406 14.5963C20.8105 14.5853 20.7411 14.5558 20.6683 14.5248L20.6682 14.5247L20.5971 14.4945C20.5026 14.4509 20.4008 14.4182 20.3718 14.4182C20.3427 14.4182 20.3063 14.4 20.2954 14.3818C20.2845 14.36 20.2482 14.3455 20.2155 14.3455C20.1864 14.3455 20.1101 14.3237 20.0483 14.2946C19.9828 14.2655 19.6703 14.1565 19.3504 14.0547C19.0306 13.9493 18.7434 13.8512 18.7107 13.8366C18.678 13.8185 18.5871 13.793 18.5108 13.7785C18.4308 13.7676 18.3545 13.7421 18.3327 13.7276C18.3145 13.7131 18.2382 13.6876 18.1691 13.6767C18.1001 13.6622 18.0092 13.6367 17.9692 13.6186C17.9292 13.604 17.8238 13.5713 17.733 13.5459C17.6941 13.5365 17.6467 13.5246 17.5988 13.5125L17.5987 13.5125L17.5986 13.5124C17.5345 13.4963 17.4698 13.4799 17.424 13.4695C17.3441 13.4477 17.2459 13.415 17.2059 13.4005C17.166 13.3823 17.0751 13.3569 17.006 13.3423C16.937 13.3314 16.857 13.306 16.8243 13.2914C16.7952 13.2769 16.6462 13.226 16.4972 13.1824C16.3481 13.1388 16.1991 13.0915 16.17 13.077C16.1373 13.0588 16.0574 13.0334 15.9883 13.0188C15.9193 13.0043 15.8466 12.9825 15.8248 12.9679C15.8066 12.9534 15.6648 12.9025 15.5158 12.8553C15.3668 12.8117 15.2287 12.7608 15.2069 12.7462C15.1887 12.7317 15.0651 12.6844 14.9343 12.6372C14.8034 12.5899 14.6835 12.5391 14.6653 12.5209C14.6435 12.5063 14.6071 12.4918 14.5781 12.4918C14.5054 12.4918 14.1637 12.3355 13.2986 11.9139C12.6916 11.6158 12.1319 11.3105 12.0083 11.2124C11.9907 11.1983 11.8501 11.1056 11.6957 11.0039L11.6812 10.9943C11.5516 10.9089 11.4316 10.8283 11.3791 10.793L11.3791 10.793L11.3541 10.7762C11.3429 10.7672 11.2733 10.7152 11.1772 10.6433L10.9761 10.4927C10.3182 10.002 9.36227 9.08246 8.95882 8.54816C8.92394 8.50165 8.89133 8.45893 8.86586 8.42556L8.86585 8.42555C8.83559 8.38591 8.81541 8.35947 8.81343 8.35552C8.80616 8.35189 8.79526 8.35916 8.78436 8.37006ZM25.0387 19.3832C25.4094 19.5068 25.5294 19.6013 25.6166 19.823C25.7038 20.0374 25.5585 20.0883 25.3804 19.9029C25.3222 19.8411 25.2168 19.7575 25.1477 19.7139C25.0787 19.6703 24.9987 19.6158 24.9733 19.5904C24.9478 19.5649 24.9042 19.5431 24.8751 19.5431C24.8279 19.5431 24.7843 19.4668 24.7843 19.3759C24.7843 19.3105 24.8352 19.3105 25.0387 19.3832ZM6.17465 10.6599C6.17465 10.6999 6.4836 11.3105 6.53812 11.376C6.55992 11.4014 6.63625 11.4995 6.70895 11.6013C6.97428 11.9575 7.38863 12.3828 7.94111 12.8553C8.07923 12.9752 8.45723 13.2587 8.57718 13.335C8.62806 13.3678 8.69349 13.4186 8.72257 13.4441C8.75164 13.4732 8.83888 13.535 8.91884 13.5822C8.9988 13.6295 9.0715 13.6803 9.0824 13.6912C9.10421 13.7203 10.1764 14.3564 10.5181 14.5454C10.7871 14.6981 11.0778 14.8434 12.0992 15.3414C12.6989 15.6358 12.8116 15.694 13.2514 15.9593C13.5821 16.1592 14.1128 16.5736 14.3781 16.8389C14.5272 16.9879 14.6544 17.1079 14.6617 17.1079C14.6726 17.1079 14.658 17.0679 14.6362 17.017C14.6144 16.9661 14.5817 16.9261 14.5635 16.9261C14.549 16.9261 14.5344 16.9043 14.5344 16.8789C14.5344 16.8498 14.5017 16.788 14.4617 16.7408C14.4327 16.7065 14.4018 16.6665 14.3815 16.6401L14.3815 16.6401L14.3636 16.6172C14.0983 16.2682 13.8802 16.0574 13.2114 15.4904C13.1205 15.4105 12.5208 15.0034 12.2409 14.8289C12.0192 14.6872 11.5722 14.4364 11.1178 14.1928C11.0088 14.131 10.9034 14.0693 10.8816 14.0547C10.8544 14.0336 10.742 13.9724 10.4071 13.7899L10.1728 13.6622C9.84569 13.4877 9.01698 12.9752 8.88249 12.8662C8.86432 12.848 8.78072 12.7935 8.70076 12.7462C8.62987 12.7043 8.55899 12.6567 8.52862 12.6363L8.5286 12.6363L8.51902 12.6299C8.50085 12.6154 8.30458 12.47 8.08286 12.3064C7.86478 12.1429 7.67578 11.9975 7.66487 11.9866C7.65397 11.972 7.56673 11.8993 7.46496 11.823C7.26142 11.6704 6.99972 11.4305 6.51267 10.9507C6.3273 10.769 6.17465 10.6381 6.17465 10.6599ZM10.3037 17.1696C10.2673 17.206 10.8779 17.773 11.2814 18.0747L11.4946 18.2334C11.5964 18.309 11.67 18.3638 11.6812 18.3727C11.7757 18.4454 12.1028 18.6344 12.1319 18.6344C12.1501 18.6344 12.1901 18.6562 12.2155 18.678C12.2409 18.7035 12.3209 18.7507 12.39 18.7871C12.8734 19.0161 13.0806 19.1069 13.1314 19.1069C13.1605 19.1069 13.1932 19.1215 13.2041 19.136C13.2332 19.1833 14.3781 19.5431 14.5054 19.5431C14.5417 19.5431 14.5999 19.5576 14.6362 19.5794C14.6689 19.5976 14.7707 19.6303 14.8616 19.6485C15.1778 19.7176 15.3522 19.7612 15.4249 19.7939C15.4649 19.8121 15.6212 19.863 15.7702 19.9102C15.9193 19.9575 16.0574 20.0047 16.0792 20.0229C16.0974 20.0374 16.1337 20.052 16.1628 20.052C16.1919 20.052 16.2646 20.0738 16.3263 20.1028C16.9224 20.3609 17.0824 20.4227 17.1042 20.4045C17.1151 20.39 17.0787 20.3427 17.0206 20.2991C16.9297 20.2264 16.7698 20.0847 16.4245 19.7612C16.37 19.7103 16.2718 19.6231 16.2064 19.5649C16.0247 19.4086 15.2614 18.8416 15.0978 18.7435C15.0179 18.6962 14.9452 18.6453 14.9343 18.6344C14.9234 18.6199 14.8906 18.5981 14.8616 18.5835C14.8325 18.569 14.7089 18.5036 14.589 18.4345C14.3636 18.3037 13.902 18.0747 13.6439 17.9656C13.564 17.9329 13.4731 17.8893 13.444 17.8748C13.4113 17.8566 13.335 17.8312 13.2696 17.8203C13.2078 17.8057 13.146 17.7766 13.1351 17.7585C13.1205 17.7403 13.0769 17.7258 13.0333 17.7258C12.9533 17.7294 12.9533 17.7294 13.026 17.7694C13.066 17.7912 13.106 17.8203 13.1169 17.8348C13.1278 17.8457 13.2078 17.9002 13.2986 17.9547C13.4913 18.0674 13.8438 18.38 13.8438 18.4382C13.8438 18.46 13.8184 18.4999 13.7893 18.5254C13.7457 18.5617 13.6403 18.5581 13.2623 18.5108C13.0006 18.4781 12.7498 18.4345 12.7062 18.4163L12.6788 18.4072L12.6788 18.4072C12.6218 18.388 12.5177 18.3529 12.4263 18.3255C11.8593 18.1547 10.718 17.5222 10.3873 17.1987C10.3545 17.1696 10.3182 17.1551 10.3037 17.1696ZM19.0451 21.6221C18.936 21.4477 17.195 19.6703 16.5153 19.0452C16.1737 18.7289 16.0828 18.6054 16.33 18.7907C16.519 18.9361 16.8243 19.1433 16.8424 19.1433C16.8497 19.1433 16.966 19.2087 17.0969 19.2887C17.2277 19.3686 17.344 19.4341 17.3586 19.4341C17.3695 19.4341 17.4567 19.4704 17.5476 19.5104C17.7838 19.6194 17.8311 19.6376 18.0964 19.7176C18.5326 19.8484 18.6925 19.8775 19.1323 19.8957C19.5612 19.9138 19.7211 19.8993 20.6407 19.7721C21.2586 19.6885 22.6507 19.6849 23.0214 19.7721C23.3122 19.8411 23.5557 19.9102 23.7047 19.972C23.7629 19.9974 23.8247 20.0156 23.8429 20.0156C23.861 20.0156 23.9046 20.0302 23.941 20.0483C23.9737 20.0665 24.1191 20.1392 24.2572 20.2083C24.5407 20.3464 24.7443 20.4554 24.7988 20.499C24.9187 20.5899 25.2313 20.7789 25.2604 20.7789C25.2786 20.7789 25.2931 20.7934 25.2931 20.8116C25.2931 20.8298 25.3258 20.8589 25.3658 20.8698C25.533 20.9207 25.413 20.9461 24.9587 20.9461C24.3335 20.9461 23.9592 21.0333 23.3485 21.3241C23.1886 21.4004 22.9778 21.4986 22.876 21.5422C22.7779 21.5822 22.6616 21.6331 22.6216 21.6512C22.3853 21.7494 22.1055 21.8075 21.7784 21.8257C21.3167 21.8475 21.3022 21.8802 21.6257 22.171C21.9383 22.4545 22.2908 22.5708 22.9851 22.6144C23.3922 22.6435 23.6575 22.7126 23.7556 22.8216C23.8247 22.9016 23.8392 23.1996 23.7883 23.4758C23.7629 23.6176 23.6938 23.6685 23.6938 23.5449C23.6938 23.4649 23.5739 23.1887 23.5048 23.1051C23.4794 23.0724 23.3922 23.007 23.3122 22.9597C23.185 22.8798 23.1268 22.8725 22.836 22.8725C22.4798 22.8761 22.3308 22.927 22.0873 23.1124C21.9201 23.2432 21.8401 23.3995 21.8401 23.5921C21.8438 23.7703 21.9674 23.952 22.2218 24.1446C22.3127 24.2173 22.3817 24.2827 22.3708 24.2936C22.3417 24.3227 21.8438 24.1882 21.8002 24.1374C21.7784 24.1083 21.742 24.0865 21.7238 24.0865C21.633 24.0865 21.3131 23.6757 21.175 23.3777C21.0623 23.1342 21.055 23.1233 20.9751 23.1124C20.9351 23.1051 20.8442 23.1305 20.7788 23.1669C20.6298 23.2469 20.1028 23.2578 19.9138 23.1814C19.6557 23.0833 19.4412 22.9779 19.4412 22.9524C19.4412 22.9343 19.5321 22.9343 19.6412 22.9452C19.8229 22.967 19.8411 22.9633 19.8411 22.9016C19.8411 22.8325 19.8193 22.8143 19.6012 22.7053C19.4667 22.6399 19.0705 22.3673 18.9215 22.24C18.5907 21.9565 18.0637 21.4368 17.9583 21.2841C17.8929 21.1896 17.951 21.1969 18.0601 21.2987C18.1546 21.3859 18.5144 21.6149 18.5617 21.6149C18.5762 21.6185 18.6525 21.6476 18.7325 21.6876C18.8124 21.7276 18.9142 21.7566 18.9542 21.7603C19.0378 21.7603 19.0887 21.6876 19.0451 21.6221Z"
      fill="white"
    />
  </svg>
);

export default Falcon;
