#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

#define M 20
#define P 2
#define N 10000
#define W 16

typedef struct {
	long id;
	double *A;
	double *B;
	double *C;
} data_t;

void multiply(const double *, const double *, double *, size_t, size_t, size_t);
void *worker(void *);

int main(int argc, char *argv[]) {
	long id;
	size_t i;
	data_t *data;
	double *A, *B, *C;
	pthread_t threads[W];

	A = malloc(sizeof(double)*M*P);
	for (i = 0; i < M*P; i++) A[i] = 42;

	B = malloc(sizeof(double)*P*N);
	for (i = 0; i < P*N; i++) B[i] = 42;

	C = malloc(sizeof(double)*M*N);
	multiply(A, B, C, M, P, N);

	for (id = 0; id < W; id++) {
		data = (data_t *)malloc(sizeof(data_t));
		data->id = id;
		data->A = A;
		data->B = B;
		data->C = C;

		if (pthread_create(&threads[id], NULL, worker, (void *)data)) {
			printf("Cannot create a thread.\n");
			exit(-1);
		}
	}

	pthread_exit(NULL);

	free((void *)A);
	free((void *)B);
	free((void *)C);
}

void *worker(void *context) {
	char bad;
	size_t i;
	double *C;
	data_t *data;

	data = (data_t *)context;

	while (1) {
		C = (double *)malloc(sizeof(double)*M*N);
		for (i = 0; i < M*N; i++) C[i] = (double)(data->id);

		multiply(data->A, data->B, C, M, P, N);

		bad = 0;

		for (i = 0; i < M*N; i++) {
			if (data->C[i] == C[i]) continue;
			printf("%6lu: %30.20e %30.20e\n", i, data->C[i], C[i]);
			bad = 1;
		}

		if (bad) {
			printf("ID: %ld\n", data->id);
			exit(-1); /* Yes, just like that. */
		}

		free((void *)C);
	}

	free((void *)data);
	pthread_exit(NULL);
}

void multiply(const double *A, const double *B, double *C, size_t m, size_t p, size_t n) {
	size_t i, j, k;
	double sum;

	for (i = 0; i < m; i++) {
		for (j = 0; j < n; j++) {
			sum = 0.0;
			for (k = 0; k < p; k++) {
				sum += A[k*m+i] * B[j*p+k];
			}
			C[j*m+i] = sum;
		}
	}
}
