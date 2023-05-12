import matplotlib.pyplot as plt
import numpy as np

def divide(arr):
    val = arr[0]
    for i in range(len(arr)):
        arr[i] /= val
        arr[i] = 1 / arr[i]
        arr[i] = round(arr[i], 2)

def generateGraphs():
    f = open("time.txt", "r")
    pipeline_small = []
    pipeline_mixture = []
    pipeline_big = []
    bsp_small = []
    bsp_mixture = []
    bsp_big = []

    threads = {1, 2, 4, 6, 8, 12}
   
    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        pipeline_small.append(time)

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        pipeline_mixture.append(time)

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        pipeline_big.append(time)

    bsp_small.append(pipeline_small[0])
    bsp_mixture.append(pipeline_mixture[0])
    bsp_big.append(pipeline_big[0])

    threads = {2, 4, 6, 8, 12}

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        bsp_small.append(time)

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        bsp_mixture.append(time)

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        bsp_big.append(time)

   
    f.close()

    divide(pipeline_small)
    divide(pipeline_mixture)
    divide(pipeline_big)
    divide(bsp_small)
    divide(bsp_mixture)
    divide(bsp_big)

    xpoints = [1, 2, 4, 6, 8, 12]
    fig1 = plt.figure("Figure 1")
    plt.title("SpeedUp Graph Pipeline")

    ypoints = np.array(pipeline_small)
    plt.plot(xpoints, ypoints, marker = 'o', label = "SMALL")

    ypoints = np.array(pipeline_mixture)
    plt.plot(xpoints, ypoints, marker = 'o', label = "MIXTURE")

    ypoints = np.array(pipeline_big)
    plt.plot(xpoints, ypoints, marker = 'o', label = "BIG")

    plt.xlabel("Number of threads")
    plt.ylabel("Speed Up")

    plt.legend()
    plt.savefig("speedup-pipeline.png")

    fig2 = plt.figure("Figure 2")
    ypoints = np.array(bsp_small)
    plt.plot(xpoints, ypoints, marker = 'o', label = "SMALL")

    ypoints = np.array(bsp_mixture)
    plt.plot(xpoints, ypoints, marker = 'o', label = "MIXTURE")

    ypoints = np.array(bsp_big)
    plt.plot(xpoints, ypoints, marker = 'o', label = "BIG")
    plt.title("SpeedUp Graph BSP")
    plt.xlabel("Number of threads")
    plt.ylabel("Speed Up")

    plt.legend()
    plt.savefig("speedup-bsp.png")


if __name__ == "__main__":
    generateGraphs()